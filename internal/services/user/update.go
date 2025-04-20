package user

import (
	"context"
	"errors"
	"fmt"

	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	"github.com/bovinxx/auth-service/internal/utils"
)

func (s *Serv) UpdateUser(ctx context.Context, id int64, oldPassword, newPassword string) error {
	user, err := s.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotExists) {
			return serviceerrs.ErrUserNotExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	ok := utils.VerifyPassword(user.Password, oldPassword)
	if !ok {
		return serviceerrs.ErrInvalidCredentials
	}

	if err := s.repo.UpdateUser(ctx, id, newPassword); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
