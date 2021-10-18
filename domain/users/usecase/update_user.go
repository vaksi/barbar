package usecase

import (
	"barbar/domain/users/entity"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *userUseCase) Update(ctx context.Context,
	id string, name, phone, email, password *string) (user entity.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "use_case.Update")
	defer span.Finish()

	oldUser, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return user, err
	}

	if oldUser == nil {
		return user, errors.New("unknown userId")
	}

	user = *oldUser

	if name != nil && name != &user.Name {
		user.Name = *name
	}

	if phone != nil && phone != &user.Phone {
		user.Phone = *phone
	}

	if email != nil && email != &user.Email {
		user.Email = *email
	}

	if password != nil && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password)) != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	user, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}
