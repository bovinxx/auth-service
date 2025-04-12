package access

import (
	"context"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/services/access/errors"
	"github.com/bovinxx/auth-service/internal/utils"
	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authPrefix = "Bearer: "
)

type userRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type sessionRepository interface {
	GetSession(ctx context.Context, id int64) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
}

type serv struct {
	userRepo    userRepository
	sessionRepo sessionRepository
	cache       cache.RedisClient
	jwtConfig   config.JWTConfig
	checker     Checker
}

func NewService(
	repo userRepository,
	sessionRepo sessionRepository,
	cache cache.RedisClient,
	jwtConfig config.JWTConfig) *serv {
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
