package user

import (
	"context"
	"fmt"

	models "github.com/bovinxx/auth-service/internal/models/user"
)

func (s *serv) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}
