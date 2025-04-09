package user

import (
	"github.com/bovinxx/auth-service/internal/client/db"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	userService "github.com/bovinxx/auth-service/internal/services/user"
)

type serv struct {
	repo      userRepo.Repository
	txManager db.TxManager
}

func NewService(repo userRepo.Repository, txManager db.TxManager) userService.Service {
	return &serv{
		repo:      repo,
		txManager: txManager,
	}
}
