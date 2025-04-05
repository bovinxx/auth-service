package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	"github.com/bovinxx/auth-service/internal/repository/user/mocks"
	"github.com/bovinxx/auth-service/internal/services/user/user"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
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
		serviceError    = fmt.Errorf("failed to delete user: %v", repositoryError)

		id = gofakeit.Int64()
	)

	tests := []struct {
		name           string
		args           args
		want           error
		err            error
		repositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx,
				id,
			},
			want: nil,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "failure case",
			args: args{
				ctx,
				id,
			},
			want: serviceError,
			err:  serviceError,
			repositoryMock: func(mc *minimock.Controller) userRepo.Repository {
				mock := mocks.NewUserRepositoryMock(t)
				mock.DeleteUserMock.Expect(ctx, id).Return(repositoryError)
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

			err := service.DeleteUser(ctx, id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, err)
		})
	}
}
