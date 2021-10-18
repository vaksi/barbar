package grpc

import (
	"barbar/domain/auth/usecase"
	"barbar/proto/auth"
	"context"
)

type AuthsServiceServer struct {
	authUseCase usecase.AuthUseCaseInterface
	auth.UnimplementedAuthServiceServer
}

func NewGrpcServer(authUseCase usecase.AuthUseCaseInterface) auth.AuthServiceServer {
	return &AuthsServiceServer{
		authUseCase: authUseCase,
	}
}

func (g AuthsServiceServer) GetAuthByEmail(ctx context.Context,
	request *auth.CheckTokenRequest) (*auth.Auth, error) {
	resp, err := g.authUseCase.CheckToken(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	return &auth.Auth{
		Uid:         resp.UID,
		UserId:      resp.UserID,
		AccessToken: resp.AccessToken,
		IsLogout:    resp.IsLogout,
	}, nil
}
