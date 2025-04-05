package user

import (
	"context"

	models "github.com/bovinxx/auth-service/internal/models/user"
)

func (s *serv) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}
