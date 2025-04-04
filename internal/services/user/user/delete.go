package user

import (
	"context"
	"fmt"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
