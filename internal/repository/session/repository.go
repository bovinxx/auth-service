package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
)

type Repository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, id int64) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	GetAllUserSessions(ctx context.Context, username string) ([]*models.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
	MarkSessionAsRevoked(ctx context.Context, refreshToken string) error
}
