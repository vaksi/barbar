package usecase

import (
	"barbar/pkg/redis"
	"context"
	"fmt"
)

func (u *authUseCase) Logout(ctx context.Context, userId string) (err error) {
	// get last login
	auths, err := u.authRepo.GetAllAuth(ctx, userId, map[string]bool{"createdAt": false}, 1, 0)
	if err != nil {
		return
	}
	if auths[0].IsLogout {
		return fmt.Errorf("user has been logout")
	}

	err = u.authRepo.DeleteDataRedis(ctx, userId)
	if err != nil && err != redis.ErrRedisKeyNil {
		return err
	}

	auths[0].IsLogout = true

	_, err = u.authRepo.Update(ctx, auths[0])
	return err
}
