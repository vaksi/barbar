package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaims struct {
	Email  string `json:"email"`
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func (cc *CustomClaims) GenerateToken(authId string) (newToken string, err error) {
	mySigningKey := []byte("AllYourBase")
	claim := CustomClaims{
		Email: cc.Email,
		UserID: cc.UserID,
		StandardClaims: jwt.StandardClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: time.Now().Add(time.Duration(1) * time.Hour).Unix(),
			Subject:   "user",
			Id:        authId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(mySigningKey)
}
