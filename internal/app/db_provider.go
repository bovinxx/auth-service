package app

import (
	"context"

	"github.com/bovinxx/auth-service/internal/client/db"
	"github.com/bovinxx/auth-service/internal/client/db/pg"
	"github.com/bovinxx/auth-service/internal/client/db/transaction"
	"github.com/bovinxx/auth-service/internal/closer"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/logger"
)

type dbProvider struct {
	pgConfig  config.PGConfig
	dbClient  db.Client
	txManager db.TxManager
}

func newDBProvider() *dbProvider {
	return &dbProvider{}
}

func (s *dbProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		db, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create db client", logger.Err(err))
		}
		closer.Add(db.Close)
		s.dbClient = db
	}

	return s.dbClient
}

func (s *dbProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatal("failed to create a new postgres config", logger.Err(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *dbProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.DBClient(ctx)
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}
