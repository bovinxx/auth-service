package access

import (
	accessService "github.com/bovinxx/auth-service/internal/services/access"
	desc "github.com/bovinxx/auth-service/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessServiceServer
	service accessService.Service
}

func NewImplementation(service accessService.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
