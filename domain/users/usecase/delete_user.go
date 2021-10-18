package usecase

import (
	"barbar/pkg/utils"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
)

func (u *userUseCase) DeleteUser(ctx context.Context, userId string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "use_case.DeleteUser")
	defer span.Finish()

	oldUser, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if oldUser == nil {
		return &utils.RequestError{Err: errors.New("invalid userId")}
	}

	err = u.userRepo.Delete(ctx, userId)

	return err
}
