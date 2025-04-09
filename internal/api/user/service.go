package user

import (
	userService "github.com/bovinxx/auth-service/internal/services/user"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserServiceServer
	service userService.Service
}

func NewImplementation(service userService.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
