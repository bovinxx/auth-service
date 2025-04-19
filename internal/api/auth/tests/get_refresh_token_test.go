package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/auth"
	mock "github.com/bovinxx/auth-service/internal/api/auth/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetRefreshToken(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserService := mock.NewMockAuthService(ctrl)
	impl := auth.NewImplementation(mockUserService)

	var (
		oldToken = gofakeit.Word()
		newToken = gofakeit.Word()
	)

	tests := []struct {
		name        string
		input       *auth_v1.GetRefreshTokenRequest
		mockSetup   func()
		wantToken   *models.Token
		wantErr     bool
		wantGRPCErr codes.Code
	}{
		{
			name: "success",
			input: &auth_v1.GetRefreshTokenRequest{
				OldRefreshToken: oldToken,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					GetRefreshToken(ctx, &models.Token{
						Token: oldToken,
					}).
					Return(&models.Token{Token: newToken}, nil)
			},
			wantToken:   &models.Token{Token: newToken},
			wantErr:     false,
			wantGRPCErr: codes.OK,
		},
		{
			name: "failure - failed to get refresh token",
			input: &auth_v1.GetRefreshTokenRequest{
				OldRefreshToken: oldToken,
			},
			mockSetup: func() {
				mockUserService.EXPECT().
					GetRefreshToken(ctx, &models.Token{
						Token: oldToken,
					}).
					Return(nil, errors.New("failed to get refresh token"))
			},
			wantToken:   nil,
			wantErr:     true,
			wantGRPCErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := impl.GetRefreshToken(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantGRPCErr, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken.Token, resp.RefreshToken)
			}
		})
	}
}
