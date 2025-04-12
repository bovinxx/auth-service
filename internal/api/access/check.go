package access

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	if _, err := s.service.Check(ctx, req.GetEndpointAddress()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
