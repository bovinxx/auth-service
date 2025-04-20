package access

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/access_v1"
)

// go generate mockgen -source=service.go -destination=service_mock/mock.go -package=mock

type Service interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type Implementation struct {
	desc.UnimplementedAccessServiceServer
	service Service
}

func NewImplementation(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
