package user

import (
	"github.com/bovinxx/auth-service/internal/client/db"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	service "github.com/bovinxx/auth-service/internal/services"
)

type serv struct {
	repo      userRepo.UserRepository
	txManager db.TxManager
}

func NewService(repo userRepo.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		repo:      repo,
		txManager: txManager,
	}
}
