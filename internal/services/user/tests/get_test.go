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

func TestService_GetUser(t *testing.T) {
	logger.InitForTests()

	type dependencies struct {
		repo *mock.MockuserRepository
	}

	type args struct {
		ctx context.Context
		id  int64
	}

	testUser := &models.User{
		ID:       1,
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 10),
		Role:     models.Role(0),
	}

	tests := []struct {
		name        string
		setup       func(d *dependencies)
		args        args
		want        *models.User
		wantErr     error
		expectError bool
	}{
		{
			name: "success",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(testUser, nil)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    testUser,
			wantErr: nil,
		},
		{
			name: "user not exists",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(nil, repoerrs.ErrUserNotExists)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:        nil,
			wantErr:     serviceerrs.ErrUserNotExists,
			expectError: true,
		},
		{
			name: "repository error",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(nil, errors.New("repository error"))
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:        nil,
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
			got, err := service.GetUser(tt.args.ctx, tt.args.id)

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
