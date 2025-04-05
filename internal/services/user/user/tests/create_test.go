package tests

import (
	"context"
	"errors"
	"testing"

	models "github.com/bovinxx/auth-service/internal/models/user"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	mocks "github.com/bovinxx/auth-service/internal/repository/user/mocks"
	"github.com/bovinxx/auth-service/internal/services/user/user"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) userRepo.Repository

	type args struct {
		ctx context.Context
		req *models.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		repositoryError = errors.New("repository error")

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, true, true, true, 10)
		isAdmin  = false

		req = &models.User{
			ID:       id,
			Name:     name,
			Email:    email,
			Password: password,
			IsAdmin:  isAdmin,
		}
	)

	tests := []struct {
		name           string
		args           args
		want           int64
		err            error
		repositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.CreateUserMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "failure case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repositoryError,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.CreateUserMock.Expect(ctx, req).Return(0, repositoryError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.repositoryMock(mc)
			service := user.NewService(userRepositoryMock, nil)

			newID, err := service.CreateUser(ctx, req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
