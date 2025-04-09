package user

import (
	"context"
	"errors"
	"fmt"

	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotExists) {
			return serviceerrs.ErrUserNotExists
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
