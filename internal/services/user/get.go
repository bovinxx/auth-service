package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bovinxx/auth-service/internal/models"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
)

func (s *Serv) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotExists) {
			return nil, serviceerrs.ErrUserNotExists
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
