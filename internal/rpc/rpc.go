package rpc

import (
	"context"
	"log/slog"

	"github.com/shakuni-dyutas/dyutas-auth/internal/app/svc/authsvc"
	"google.golang.org/grpc"
)

type RpcConfig struct {
	Logger *slog.Logger
}

func InitUserServiceWith(server grpc.ServiceRegistrar, cnf RpcConfig, authSvc authsvc.AuthService) {
	if server == nil {
		panic("grpc server registrar isn't configured while initializing UserService")
	}
	if authSvc == nil {
		panic("auth service isn't configured while initializing UserService")
	}

	RegisterUserServiceServer(server, &UserServiceServerImpl{
		logger:  cnf.Logger,
		authSvc: authSvc,
	})
}

type UserServiceServerImpl struct {
	UnimplementedUserServiceServer

	logger  *slog.Logger
	authSvc authsvc.AuthService
}

func (s *UserServiceServerImpl) GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	_, err := s.authSvc.GetSelfInfo(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &GetUserInfoResponse{}, nil
}
