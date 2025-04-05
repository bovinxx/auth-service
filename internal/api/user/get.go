package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.service.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return converter.ToGetResponseFromUser(user), nil
}
