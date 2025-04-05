package tests

import (
	"context"
	"errors"
	"testing"

	api "github.com/bovinxx/auth-service/internal/api/user"
	models "github.com/bovinxx/auth-service/internal/models/user"
	service "github.com/bovinxx/auth-service/internal/services/user"
	mocks "github.com/bovinxx/auth-service/internal/services/user/mocks"
	"github.com/bovinxx/auth-service/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceError = errors.New("service error")

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, true, true, true, 10)
		isAdmin  = false

		req = &user_v1.CreateRequest{
			Name:     name,
			Email:    email,
			Password: password,
			IsAdmin:  isAdmin,
		}
	)

	tests := []struct {
		name        string
		args        args
		want        int64
		err         error
		serviceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			serviceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(t)
				mock.CreateUserMock.Expect(ctx, &models.User{
					ID:       id,
					Name:     name,
					Password: password,
					IsAdmin:  isAdmin,
				}).Return(id, nil)
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
			err:  serviceError,
			serviceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(t)
				mock.CreateUserMock.Expect(ctx, nil).Return(0, serviceError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.serviceMock(mc)
			api := api.NewImplementation(userServiceMock)

			newID, err := api.Create(ctx, req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
