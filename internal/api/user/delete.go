package user

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
