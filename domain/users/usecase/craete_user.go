package usecase

import (
	"barbar/domain/users/entity"
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *userUseCase) Register(ctx context.Context, name, phone, email, password string) (user entity.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "use_case.Create")
	defer span.Finish()

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user, err = u.userRepo.Create(ctx, entity.User{
		UID:       uuid.New().String(),
		Name:      name,
		Phone:     phone,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return user, err
	}

	return user, nil
}
