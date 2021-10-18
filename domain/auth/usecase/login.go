package usecase

import (
	"barbar/domain/auth/entity"
	"barbar/domain/auth/infrastructure"
	"barbar/pkg/jwt"
	"barbar/pkg/utils"
	"barbar/proto/users"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *authUseCase) Login(ctx context.Context, email, password string) (auth entity.Auth, err error) {
	userSvc := infrastructure.NewUserService()

	user, err := userSvc.GetUserByEmail(ctx, &users.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return auth, err
	} else if user == nil {
		return auth, &utils.RequestError{
			Err: errors.New("invalid email"),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return auth, &utils.RequestError{
			Err: errors.New("invalid password"),
		}
	}

	authId := uuid.New().String()

	token := jwt.CustomClaims{
		Email:  user.Email,
		UserID: user.Uid,
	}
	newToken, err := token.GenerateToken(authId)
	if err != nil {
		return auth, err
	}

	auth = entity.Auth{
		UID:         authId,
		UserID:      user.Uid,
		AccessToken: newToken,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsLogout:    false,
	}

	// set to redis
	err = u.authRepo.CreateToRedis(ctx, auth)
	if err != nil {
		return auth, err
	}

	auth, err = u.authRepo.Create(ctx, auth)

	return auth, err
}
