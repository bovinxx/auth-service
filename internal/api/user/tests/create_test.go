package user

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/user"
	mock "github.com/bovinxx/auth-service/internal/api/user/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/internal/models"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	"github.com/bovinxx/auth-service/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserService := mock.NewMockUserService(ctrl)
	impl := user.NewImplementation(mockUserService)

	var (
		id       = gofakeit.Int64()
		username = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = models.Role(0)
	)

	tests := []struct {
		name        string
		input       *user_v1.CreateRequest
		mockSetup   func()
		wantID      int64
		wantErr     bool
		wantGRPCErr codes.Code
	}{
		{
			name: "success",
			input: &user_v1.CreateRequest{
				Name:     username,
				Email:    email,
				Password: password,
				Role:     user_v1.Role(role),
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					CreateUser(ctx, &models.User{
						Name:     username,
						Email:    email,
						Password: password,
						Role:     role,
					}).
					Return(id, nil)
			},
			wantID:      id,
			wantErr:     false,
			wantGRPCErr: codes.OK,
		},
		{
			name: "failure - already exists",
			input: &user_v1.CreateRequest{
				Name:     username,
				Email:    email,
				Password: password,
				Role:     user_v1.Role(role),
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					CreateUser(ctx, &models.User{
						Name:     username,
						Email:    email,
						Password: password,
						Role:     role,
					}).
					Return(int64(0), serviceerrs.ErrUserAlreadyExists)
			},
			wantID:      0,
			wantErr:     true,
			wantGRPCErr: codes.AlreadyExists,
		},
		{
			name: "failure - internal error",
			input: &user_v1.CreateRequest{
				Name:     username,
				Email:    email,
				Password: password,
				Role:     user_v1.Role(role),
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					CreateUser(ctx, &models.User{
						Name:     username,
						Email:    email,
						Password: password,
						Role:     role,
					}).
					Return(int64(0), errors.New("internal error"))
			},
			wantID:      0,
			wantErr:     true,
			wantGRPCErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := impl.Create(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error state: got error = %v, wantErr = %v", err, tt.wantErr)
			}

			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got: %v", err)
				}
				if st.Code() != tt.wantGRPCErr {
					t.Errorf("unexpected gRPC error code: got = %v, want = %v", st.Code(), tt.wantGRPCErr)
				}
			}

			if resp != nil && resp.Id != tt.wantID {
				t.Errorf("unexpected ID: got = %d, want = %d", resp.Id, tt.wantID)
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("expected non-nil response, got nil")
			}
		})
	}
}
