package access

import "sync"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Rule struct {
	AllowedRoles []Role
	IsPublic     bool
}

type Checker interface {
	HasAccess(userRole Role, endpoint string) bool
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

func (c *StaticChecker) HasAccess(userRole Role, endpoint string) bool {
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
