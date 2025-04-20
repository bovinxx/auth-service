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
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdate(t *testing.T) {
	logger.InitForTests()

	type mockRepo struct {
		userRepository *mock.MockuserRepository
	}

	type args struct {
		ctx         context.Context
		id          int64
		oldPassword string
		newPassword string
	}

	tests := []struct {
		name        string
		prepare     func(m *mockRepo)
		args        args
		wantErr     error
		expectError bool
	}{
		{
			name: "successful password update",
			prepare: func(m *mockRepo) {
				hashedPassword, _ := utils.HashPassword("correctPassword")
				testUser := &models.User{
					ID:       1,
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: string(hashedPassword),
					Role:     models.RoleUser,
				}

				m.userRepository.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(testUser, nil)
				m.userRepository.EXPECT().
					UpdateUser(gomock.Any(), int64(1), gomock.Any()).
					Return(nil)
			},
			args: args{
				ctx:         context.Background(),
				id:          1,
				oldPassword: "correctPassword",
				newPassword: "newSecurePassword123",
			},
			wantErr:     nil,
			expectError: false,
		},
		{
			name: "user not found",
			prepare: func(m *mockRepo) {
				m.userRepository.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(nil, repoerrs.ErrUserNotExists)
			},
			args: args{
				ctx:         context.Background(),
				id:          1,
				oldPassword: "anyPassword",
				newPassword: "newPassword",
			},
			wantErr:     serviceerrs.ErrUserNotExists,
			expectError: true,
		},
		{
			name: "incorrect old password",
			prepare: func(m *mockRepo) {
				hashedPassword, _ := utils.HashPassword("correctPassword")
				testUser := &models.User{
					ID:       1,
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: string(hashedPassword),
					Role:     models.RoleUser,
				}

				m.userRepository.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(testUser, nil)
			},
			args: args{
				ctx:         context.Background(),
				id:          1,
				oldPassword: "wrongPassword",
				newPassword: "newPassword",
			},
			wantErr:     serviceerrs.ErrInvalidCredentials,
			expectError: true,
		},
		{
			name: "repository error during update",
			prepare: func(m *mockRepo) {
				hashedPassword, _ := utils.HashPassword("correctPassword")
				testUser := &models.User{
					ID:       1,
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Password: string(hashedPassword),
					Role:     models.RoleUser,
				}

				m.userRepository.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(testUser, nil)
				m.userRepository.EXPECT().
					UpdateUser(gomock.Any(), int64(1), gomock.Any()).
					Return(errors.New("database connection failed"))
			},
			args: args{
				ctx:         context.Background(),
				id:          1,
				oldPassword: "correctPassword",
				newPassword: "newPassword",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := &mockRepo{
				userRepository: mock.NewMockuserRepository(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(mocks)
			}

			service := user.NewService(mocks.userRepository)
			err := service.UpdateUser(
				tt.args.ctx,
				tt.args.id,
				tt.args.oldPassword,
				tt.args.newPassword,
			)

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
