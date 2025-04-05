package user

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateUser(
		ctx,
		req.GetId(),
		req.GetName().Value,
		req.GetOldPassword().Value,
		req.GetNewPassword().Value,
	)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
