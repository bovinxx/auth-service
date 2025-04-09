package auth

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bovinxx/auth-service/internal/client/db"
	"github.com/bovinxx/auth-service/internal/models"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	"github.com/bovinxx/auth-service/internal/repository/user/errors"
	"github.com/bovinxx/auth-service/internal/repository/user/user/converter"
	"github.com/bovinxx/auth-service/internal/repository/user/user/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableName = "users"

	idColumn       = "id"
	nameColumn     = "name"
	emailColumn    = "email"
	passwordColumn = "password"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) (userRepo.Repository, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	if _, err := r.GetUserByUsername(ctx, user.Name); err == nil {
		return 0, errors.ErrUserAlreadyExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %v", err)
	}

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.Name, user.Email, hashPassword, user.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "userRepository.CreateUser",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	return id, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	q := db.Query{
		Name:     "userRepository.GetUserByID",
		QueryRaw: query,
	}

	user := &model.User{}
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	return converter.ToUserFromRepo(user), nil
}

func (r *repo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: username})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	q := db.Query{
		Name:     "userRepository.GetUserByUsername",
		QueryRaw: query,
	}

	user := &model.User{}
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		return nil, errors.ErrUserNotExists
	}

	return converter.ToUserFromRepo(user), nil
}

func (r *repo) UpdateUser(ctx context.Context, id int64, newPassword string) error {
	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	updateBuilder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(passwordColumn, hashNewPassword).
		Where(sq.Eq{idColumn: id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	q := db.Query{
		Name:     "userRepository.UpdateUser",
		QueryRaw: query,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	deleteBuilder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %v", err)
	}

	q := db.Query{
		Name:     "userRepository.DeleteUser",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrUserNotExists
	}

	return nil
}
