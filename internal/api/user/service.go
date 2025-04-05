package user

import (
	UserService "github.com/bovinxx/auth-service/internal/services"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserServiceServer
	service UserService.UserService
}

func NewImplementation(service UserService.UserService) *Implementation {
	return &Implementation{
		service: service,
	}
}
