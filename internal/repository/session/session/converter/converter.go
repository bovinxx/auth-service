package converter

import (
	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/repository/session/session/model"
)

func ToSessionFromRepo(session *model.Session) *models.Session {
	return &models.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		CreatedAt:    session.CreatedAt,
		ExpiresAt:    session.ExpiresAt,
		RevokedAt:    session.RevokedAt,
	}
}
