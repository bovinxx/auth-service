package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models "github.com/bovinxx/auth-service/internal/models/user"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	"github.com/bovinxx/auth-service/internal/repository/user/mocks"
	"github.com/bovinxx/auth-service/internal/services/user/user"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) userRepo.Repository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		repositoryError = errors.New("repository error")
		serviceError    = fmt.Errorf("failed to get user: %v", repositoryError)

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
		want           *models.User
		err            error
		repositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx,
				id,
			},
			want: &models.User{
				ID:       id,
				Name:     name,
				Email:    email,
				Password: password,
				IsAdmin:  isAdmin,
			},
			err: nil,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.GetUserByIDMock.Expect(ctx, req.ID).Return(&models.User{
					ID:       id,
					Name:     name,
					Email:    email,
					Password: password,
					IsAdmin:  isAdmin,
				}, nil)
				return mock
			},
		},
		{
			name: "failure case",
			args: args{
				ctx,
				id,
			},
			want: nil,
			err:  serviceError,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.GetUserByIDMock.Expect(ctx, req.ID).Return(nil, repositoryError)
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

			newID, err := service.GetUser(ctx, req.ID)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
