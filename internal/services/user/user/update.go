package user

import (
	"context"
	"fmt"
)

func (s *serv) UpdateUser(ctx context.Context, id int64, name, oldPassword, newPassword string) error {
	if err := s.repo.UpdateUser(ctx, id, name, oldPassword, newPassword); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
