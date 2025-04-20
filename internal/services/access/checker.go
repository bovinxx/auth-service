package access

import (
	"context"
	"strings"
	"sync"

	"github.com/bovinxx/auth-service/internal/models"
	serverrs "github.com/bovinxx/auth-service/internal/services/access/errors"
)

type Checker interface {
	HasAccess(userRole models.Role, endpoint string) bool
}

type StaticChecker struct {
	accessRules   map[string]Rule
	accessRulesMu sync.RWMutex
}

func NewStaticChecker(rules map[string]Rule) Checker {
	return &StaticChecker{
		accessRules:   rules,
		accessRulesMu: sync.RWMutex{},
	}
}

func (c *StaticChecker) HasAccess(userRole models.Role, endpoint string) bool {
	c.accessRulesMu.Lock()
	defer c.accessRulesMu.Unlock()

	rule, ok := c.accessRules[endpoint]
	if !ok {
		return false
	}

	if rule.IsPublic {
		return true
	}

	for _, role := range rule.AllowedRoles {
		if role == userRole {
			return true
		}
	}

	return false
}

func (s *Serv) Check(ctx context.Context, endpoint string) (bool, error) {
	endpoint = strings.ToLower(endpoint)
	if s.isEndpointPublic(endpoint) {
		return true, nil
	}

	accessToken, err := s.extractBearerToken(ctx)
	if err != nil {
		return false, err
	}

	claims, err := s.verifyToken(accessToken)
	if err != nil {
		return false, err
	}

	if ok := s.checker.HasAccess(models.Role(claims.Role), endpoint); ok {
		return false, nil
	}

	return true, serverrs.ErrAccessDenied
}

func (s *Serv) isEndpointPublic(endpoint string) bool {
	rule, ok := s.accessConfig.AccessRule()[endpoint]
	return ok && rule.Public
}
