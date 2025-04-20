package access

import (
	"context"
	"time"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/models"
)

type Rule struct {
	AllowedRoles []models.Role
	IsPublic     bool
}

const (
	RoleAdmin             models.Role = models.RoleAdmin
	RoleUser              models.Role = models.RoleUser
	authPrefix                        = "Bearer: "
	sessionCacheKeyPrefix             = "auth:session:sessionID"
	cacheExpTime                      = 10 * time.Minute
)

type userRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type sessionRepository interface {
	GetSession(ctx context.Context, id int64) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
}

type Serv struct {
	userRepo     userRepository
	sessionRepo  sessionRepository
	cache        cache.RedisClient
	jwtConfig    config.JWTConfig
	accessConfig config.AccessConfig
	checker      Checker
}

func NewService(
	repo userRepository,
	sessionRepo sessionRepository,
	cache cache.RedisClient,
	jwtConfig config.JWTConfig,
	accessConfig config.AccessConfig) *Serv {
	return &Serv{
		userRepo:     repo,
		sessionRepo:  sessionRepo,
		cache:        cache,
		jwtConfig:    jwtConfig,
		accessConfig: accessConfig,
		checker:      NewStaticChecker(buildRulesFromConfig(accessConfig.AccessRule())),
	}
}

func buildRulesFromConfig(cfg map[string]config.AccessRuleConfig) map[string]Rule {
	rules := make(map[string]Rule)

	for endpoint, rule := range cfg {
		roles := make([]models.Role, 0, len(rule.Role))
		for _, r := range rule.Role {
			roles = append(roles, models.Role(r))
		}

		rules[endpoint] = Rule{
			AllowedRoles: roles,
			IsPublic:     rule.Public,
		}
	}

	return rules
}
