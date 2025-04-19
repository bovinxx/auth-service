package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/auth"
	mock "github.com/bovinxx/auth-service/internal/api/auth/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	"github.com/bovinxx/auth-service/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserService := mock.NewMockAuthService(ctrl)
	impl := auth.NewImplementation(mockUserService)

	username := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, false, 10)
	token := gofakeit.UUID()

	tests := []struct {
		name        string
		input       *auth_v1.LoginRequest
		mockSetup   func()
		wantToken   string
		wantErr     bool
		wantGRPCErr codes.Code
	}{
		{
			name: "success",
			input: &auth_v1.LoginRequest{
				Username: username,
				Password: password,
				Role:     auth_v1.Role_ROLE_USER,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Login(ctx, gomock.Any()).
					Return(token, nil)
			},
			wantToken:   token,
			wantErr:     false,
			wantGRPCErr: codes.OK,
		},
		{
			name: "failure - invalid credentials",
			input: &auth_v1.LoginRequest{
				Username: username,
				Password: password,
				Role:     auth_v1.Role_ROLE_USER,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Login(ctx, gomock.Any()).
					Return("", serviceerrs.ErrInvalidCredentials)
			},
			wantErr:     true,
			wantGRPCErr: codes.Aborted,
		},
		{
			name: "failure - access denied",
			input: &auth_v1.LoginRequest{
				Username: username,
				Password: password,
				Role:     auth_v1.Role_ROLE_USER,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Login(ctx, gomock.Any()).
					Return("", serviceerrs.ErrAccessDenied)
			},
			wantErr:     true,
			wantGRPCErr: codes.Aborted,
		},
		{
			name: "failure - internal error",
			input: &auth_v1.LoginRequest{
				Username: username,
				Password: password,
				Role:     auth_v1.Role_ROLE_USER,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Login(ctx, gomock.Any()).
					Return("", errors.New("unexpected error"))
			},
			wantErr:     true,
			wantGRPCErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := impl.Login(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantGRPCErr, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantToken, resp.RefreshToken)
			}
		})
	}
}
