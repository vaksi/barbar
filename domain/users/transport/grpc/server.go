package grpc

import (
	"barbar/domain/users/usecase"
	"barbar/proto/users"
	"context"
)

type UsersServiceServer struct {
	userUseCase usecase.UserUseCaseInterface
	users.UnimplementedUsersServiceServer
}

func NewGrpcServer(userUseCase usecase.UserUseCaseInterface) users.UsersServiceServer {
	return &UsersServiceServer{
		userUseCase: userUseCase,
	}
}

func (g UsersServiceServer) GetUserByEmail(ctx context.Context,
	request *users.GetUserByEmailRequest) (*users.User, error) {
	resp, err := g.userUseCase.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	return &users.User{
		Uid:      resp.UID,
		Name:     resp.Name,
		Phone:    resp.Phone,
		Email:    resp.Email,
		Password: resp.Password,
	}, nil
}
