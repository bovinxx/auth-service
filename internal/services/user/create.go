package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bovinxx/auth-service/internal/models"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
)

func (s *Serv) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserAlreadyExists) {
			return 0, serviceerrs.ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("failed to create a new user: %w", err)
	}
	return id, nil
}
