package access

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/access_v1"
)

type AccessService interface {
	Check(ctx context.Context, endpoint string) (bool, error)
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
