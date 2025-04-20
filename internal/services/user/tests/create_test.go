package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/internal/models"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	"github.com/bovinxx/auth-service/internal/services/user"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	mock "github.com/bovinxx/auth-service/internal/services/user/service_mock"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	logger.InitForTests()

	type dependencies struct {
		repo *mock.MockuserRepository
	}

	type args struct {
		ctx  context.Context
		user *models.User
	}

	tests := []struct {
		name        string
		setup       func(d *dependencies)
		args        args
		want        int64
		wantErr     error
		expectError bool
	}{
		{
			name: "success",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
			},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 10),
					Role:     models.Role(0),
				},
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "user already exists",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(int64(0), repoerrs.ErrUserAlreadyExists)
			},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 10),
					Role:     models.Role(0),
				},
			},
			want:        0,
			wantErr:     serviceerrs.ErrUserAlreadyExists,
			expectError: true,
		},
		{
			name: "repository error",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("repository error"))
			},
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 10),
					Role:     models.Role(0),
				},
			},
			want:        0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := &dependencies{
				repo: mock.NewMockuserRepository(ctrl),
			}

			if tt.setup != nil {
				tt.setup(d)
			}

			service := user.NewService(d.repo)
			got, err := service.CreateUser(tt.args.ctx, tt.args.user)

			if tt.expectError {
				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
