package access

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccessService interface {
	Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error)
}

type Implementation struct {
	desc.UnimplementedAccessServiceServer
	service AccessService
}

func NewImplementation(service AccessService) *Implementation {
	return &Implementation{
		service: service,
	}
}
