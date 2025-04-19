package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/auth"
	mock "github.com/bovinxx/auth-service/internal/api/auth/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogout(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserService := mock.NewMockAuthService(ctrl)
	impl := auth.NewImplementation(mockUserService)

	token := gofakeit.UUID()

	tests := []struct {
		name        string
		input       *auth_v1.LogoutRequest
		mockSetup   func()
		wantToken   string
		wantErr     bool
		wantGRPCErr codes.Code
	}{
		{
			name: "success",
			input: &auth_v1.LogoutRequest{
				RefreshToken: token,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Logout(ctx, gomock.Any()).
					Return(nil)
			},
			wantToken:   token,
			wantErr:     false,
			wantGRPCErr: codes.OK,
		},
		{
			name: "failure - internal error",
			input: &auth_v1.LogoutRequest{
				RefreshToken: token,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					Logout(ctx, gomock.Any()).
					Return(errors.New("internal error"))
			},
			wantErr:     true,
			wantGRPCErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := impl.Logout(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantGRPCErr, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
