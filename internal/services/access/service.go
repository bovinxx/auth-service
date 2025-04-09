package service

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error)
}
