package access

import (
	"context"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	sessionRepo "github.com/bovinxx/auth-service/internal/repository/session"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	accessService "github.com/bovinxx/auth-service/internal/services/access"
	"github.com/bovinxx/auth-service/internal/services/access/errors"
	"github.com/bovinxx/auth-service/internal/utils"
	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authPrefix = "Bearer: "
)

type serv struct {
	userRepo    userRepo.Repository
	sessionRepo sessionRepo.Repository
	cache       cache.RedisClient
	jwtConfig   config.JWTConfig
	checker     Checker
}

func NewService(
	repo userRepo.Repository,
	sessionRepo sessionRepo.Repository,
	cache cache.RedisClient,
	jwtConfig config.JWTConfig) accessService.Service {
	return &serv{
		userRepo:    repo,
		sessionRepo: sessionRepo,
		cache:       cache,
		jwtConfig:   jwtConfig,
		checker:     NewStaticChecker(nil),
	}
}

func (s *serv) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	accessToken, err := extractBearerToken(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := utils.VerifyToken(accessToken, []byte(s.jwtConfig.AccessTokenSecret()))
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if ok := s.checker.HasAccess(Role(claims.Role), req.GetEndpointAddress()); ok {
		return nil, nil
	}

	return nil, errors.ErrAccessDenied
}
