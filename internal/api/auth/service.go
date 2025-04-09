package auth

import (
	authService "github.com/bovinxx/auth-service/internal/services/auth"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthServiceServer
	service authService.Service
}

func NewImplementation(service authService.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
