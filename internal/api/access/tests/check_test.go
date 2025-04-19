package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bovinxx/auth-service/internal/api/access"
	mock "github.com/bovinxx/auth-service/internal/api/access/service_mock"
	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/access/errors"
	"github.com/bovinxx/auth-service/pkg/access_v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCheck(t *testing.T) {
	logger.InitForTests()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockAccessService := mock.NewMockAccessService(ctrl)
	impl := access.NewImplementation(mockAccessService)

	endpoint := "/endpoint"

	tests := []struct {
		name        string
		err         error
		wantGRPCErr codes.Code
	}{
		{"access denied", serviceerrs.ErrAccessDenied, codes.PermissionDenied},
		{"no auth header", serviceerrs.ErrNoAuthHeader, codes.Unauthenticated},
		{"invalid token", serviceerrs.ErrInvalidToken, codes.Unauthenticated},
		{"invalid auth", serviceerrs.ErrInvalidAuth, codes.PermissionDenied},
		{"no metadata", serviceerrs.ErrNoMetadata, codes.FailedPrecondition},
		{"internal error", errors.New("some internal error"), codes.Internal},
		{"success", nil, codes.OK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				mockAccessService.EXPECT().Check(ctx, endpoint).Return(false, tt.err)
			} else {
				mockAccessService.EXPECT().Check(ctx, endpoint).Return(true, nil)
			}

			req := &access_v1.CheckRequest{EndpointAddress: endpoint}
			resp, err := impl.Check(ctx, req)

			if tt.wantGRPCErr == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantGRPCErr, st.Code())
				assert.Nil(t, resp)
			}
		})
	}
}
