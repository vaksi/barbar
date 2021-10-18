package usecase

import (
	"barbar/domain/auth/entity"
	"barbar/domain/auth/repository"
	"context"
)

type AuthUseCaseInterface interface {
	Login(ctx context.Context, email, password string) (entity.Auth, error)
	Logout(ctx context.Context, userId string) (err error)
}

type authUseCase struct {
	authRepo repository.AuthRepositoryInterface
}

func NewAuthUseCase(userRepository repository.AuthRepositoryInterface) *authUseCase {
	return &authUseCase{
		authRepo: userRepository,
	}
}
