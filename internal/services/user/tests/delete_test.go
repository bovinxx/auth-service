package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/logger"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/user/errors"
	"github.com/bovinxx/auth-service/internal/services/user"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	mock "github.com/bovinxx/auth-service/internal/services/user/service_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDelete(t *testing.T) {
	logger.InitForTests()

	type dependencies struct {
		repo *mock.MockuserRepository
	}

	type args struct {
		ctx context.Context
		id  int64
	}

	tests := []struct {
		name        string
		setup       func(d *dependencies)
		args        args
		wantErr     error
		expectError bool
	}{
		{
			name: "success",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					DeleteUser(gomock.Any(), int64(1)).
					Return(nil)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr:     nil,
			expectError: false,
		},
		{
			name: "user not exists",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					DeleteUser(gomock.Any(), int64(1)).
					Return(repoerrs.ErrUserNotExists)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr:     serviceerrs.ErrUserNotExists,
			expectError: true,
		},
		{
			name: "repository error",
			setup: func(d *dependencies) {
				d.repo.EXPECT().
					DeleteUser(gomock.Any(), int64(1)).
					Return(errors.New("repository error"))
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
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
			err := service.DeleteUser(tt.args.ctx, tt.args.id)

			if tt.expectError {
				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				return
			}

			require.NoError(t, err)
		})
	}
}
