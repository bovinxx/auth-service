package auth

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bovinxx/auth-service/internal/client/db"
	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/repository/session/converter"
	"github.com/bovinxx/auth-service/internal/repository/session/errors"
	"github.com/bovinxx/auth-service/internal/repository/session/model"
)

const (
	tableName = "sessions"

	idColumn           = "id"
	userIDColumn       = "user_id"
	refreshTokenColumn = "refresh_token"
	createdAtColumn    = "created_at"
	expiresAtColumn    = "expires_at"
	revokedAtColumn    = "revoked_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) (*repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) CreateSession(ctx context.Context, session *models.Session) error {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, refreshTokenColumn, createdAtColumn, expiresAtColumn, revokedAtColumn).
		Values(session.UserID, session.RefreshToken, session.CreatedAt, session.ExpiresAt, session.RevokedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "sessionRepository.CreateSession",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert session: %v", err)
	}

	return nil
}

func (r *repo) GetSession(ctx context.Context, id int64) (*models.Session, error) {
	builderSelect := sq.Select(userIDColumn, refreshTokenColumn, createdAtColumn, expiresAtColumn, revokedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	q := db.Query{
		Name:     "sessionRepository.GetSession",
		QueryRaw: query,
	}

	session := &model.Session{}
	err = r.db.DB().ScanOneContext(ctx, session, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select session: %v", err)
	}

	return converter.ToSessionFromRepo(session), nil
}

func (r *repo) GetAllUserSessions(ctx context.Context, username string) ([]*models.Session, error) {
	return nil, nil
}

func (r *repo) GetSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	builderSelect := sq.Select(userIDColumn, refreshTokenColumn, createdAtColumn, expiresAtColumn, revokedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{refreshTokenColumn: token})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	q := db.Query{
		Name:     "sessionRepository.GetSessionByToken",
		QueryRaw: query,
	}

	session := &model.Session{}
	err = r.db.DB().ScanOneContext(ctx, session, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select session: %v", err)
	}

	return converter.ToSessionFromRepo(session), nil
}

func (r *repo) DeleteSession(ctx context.Context, refreshToken string) error {
	deleteBuilder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{refreshTokenColumn: refreshToken})

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %v", err)
	}

	q := db.Query{
		Name:     "sessionRepository.DeleteSesssion",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrSessionNotExists
	}

	return nil
}

func (r *repo) MarkSessionAsRevoked(ctx context.Context, refreshToken string) error {
	now := time.Now()
	updateBuilder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(revokedAtColumn, &now).
		Where(sq.Eq{refreshTokenColumn: refreshToken})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build mark session as revoked: %v", err)
	}

	q := db.Query{
		Name:     "sessionRepository.MarkSessionAsRevoked",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to mark session as revoked: %v", err)
	}

	return nil
}
