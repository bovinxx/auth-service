package user

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/user"
	mock "github.com/bovinxx/auth-service/internal/api/user/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	"github.com/bovinxx/auth-service/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdate(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserService := mock.NewMockUserService(ctrl)
	impl := user.NewImplementation(mockUserService)

	var (
		id          = gofakeit.Int64()
		oldPassword = gofakeit.Password(true, true, true, true, false, 10)
		newPassword = gofakeit.Password(true, true, true, true, false, 10)
	)

	tests := []struct {
		name        string
		input       *user_v1.UpdateRequest
		mockSetup   func()
		wantID      int64
		wantErr     bool
		wantGRPCErr codes.Code
	}{
		{
			name: "success",
			input: &user_v1.UpdateRequest{
				Id:          id,
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					UpdateUser(ctx, id, oldPassword, newPassword).
					Return(nil)
			},
			wantID:      id,
			wantErr:     false,
			wantGRPCErr: codes.OK,
		},
		{
			name: "failure - user not exists",
			input: &user_v1.UpdateRequest{
				Id:          id,
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					UpdateUser(ctx, id, oldPassword, newPassword).
					Return(serviceerrs.ErrUserNotExists)
			},
			wantID:      0,
			wantErr:     true,
			wantGRPCErr: codes.NotFound,
		},
		{
			name: "failure - invalid credentials",
			input: &user_v1.UpdateRequest{
				Id:          id,
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					UpdateUser(ctx, id, oldPassword, newPassword).
					Return(serviceerrs.ErrInvalidCredentials)
			},
			wantID:      0,
			wantErr:     true,
			wantGRPCErr: codes.Unauthenticated,
		},
		{
			name: "failure - internal error",
			input: &user_v1.UpdateRequest{
				Id:          id,
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					UpdateUser(ctx, id, oldPassword, newPassword).
					Return(errors.New("internal error"))
			},
			wantID:      0,
			wantErr:     true,
			wantGRPCErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := impl.Update(ctx, tt.input)

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
		})
	}
}
