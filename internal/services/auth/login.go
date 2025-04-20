package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	serverrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	"github.com/bovinxx/auth-service/internal/utils"
)

func (s *Serv) Login(ctx context.Context, req *models.User) (string, error) {
	user, err := s.getUserByUsername(ctx, req.Name)
	if err != nil {
		return "", err
	}

	isCorrectPassword := utils.VerifyPassword(user.Password, req.Password)
	isCorrectRole := req.Role == user.Role

	if !isCorrectPassword {
		return "", serverrs.ErrInvalidCredentials
	}

	if !isCorrectRole {
		return "", serverrs.ErrAccessDenied
	}

	refreshToken, err := s.createRefreshToken(user.ID, user.Name, user.Role)
	if err != nil {
		return "", err
	}

	err = s.createSession(ctx, user.ID, refreshToken)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
