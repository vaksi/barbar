package usecase

import (
	"barbar/domain/users/entity"
	"barbar/pkg/utils"
	"context"
	"errors"
)

func (u *userUseCase) GetById(ctx context.Context, id string) (user *entity.User, err error) {
	user, err = u.userRepo.GetById(ctx, id)
	if err != nil {
		return user, err
	}
	if user == nil {
		return user, &utils.RequestError{
			Err: errors.New("invalid userId"),
		}
	}
	return user, nil
}

func (u *userUseCase) GetAllUser(ctx context.Context,
	search string, sortBy map[string]bool, limit, offset int64) (users []entity.User, err error) {
	users, err = u.userRepo.GetAllUser(ctx, search, sortBy, limit, offset)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *userUseCase) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	user, err = u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	return user, nil
}
