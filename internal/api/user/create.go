package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.service.CreateUser(ctx, converter.ToUserFromCreateRequest(req))
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
