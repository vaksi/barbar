package usecase

import (
	"barbar/domain/users/entity"
	"barbar/domain/users/repository"
	"context"
)

type UserUseCaseInterface interface {
	Register(ctx context.Context, name, phone, email, password string) (entity.User, error)
	Update(ctx context.Context,
		id string, name, phone, email, password *string) (entity.User, error)
	Delete(ctx context.Context, id string) error

	GetById(ctx context.Context, id string) (*entity.User, error)
	GetAllUser(ctx context.Context, search string, sortBy map[string]bool, limit, offset int64) ([]entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserUseCase(userRepository repository.UserRepositoryInterface) *userUseCase {
	return &userUseCase{
		userRepo: userRepository,
	}
}
