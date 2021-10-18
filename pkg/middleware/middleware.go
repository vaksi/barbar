package middleware

import (
	"barbar/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Authorization struct {
	cache *redis.Redis
}

func NewAuthMiddleware(cache *redis.Redis) *Authorization {
	return &Authorization{
		cache: cache,
	}
}

const userInfoKey = "userInfo"

func (a *Authorization) CheckToken(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user-service/api/v1/users" && r.Method == "POST" {
			handler.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/auth-service/api/v1/login" {
			handler.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return []byte("AllYourBase"), nil
		})
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "invalid claims", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), userInfoKey, claims)
		r = r.WithContext(ctx)

		// check validity internal token
		var tokenFromRedis string

		user := GetUserInfoFromContext(ctx)

		if user == nil {
			http.Error(w, "invalid claims because error parsing claims", http.StatusBadRequest)
			return
		}
		fmt.Println(user)
		err = a.cache.GetDataRedis(r.Context(), "users-"+user.UserID, &tokenFromRedis)
		if err != nil {
			logger.Error(err)
			http.Error(w, "invalid claims because redis err = "+err.Error(), http.StatusBadRequest)
			return
		}
		if tokenFromRedis != tokenString {
			http.Error(w, "invalid claims because unknown token", http.StatusBadRequest)
		}
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type userInfo struct {
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	AuthID string `json:"jti"`
	Sub    string `json:"sub"`
	UserID string `json:"userId"`
}

func GetUserInfoFromContext(ctx context.Context) *userInfo {
	var usr *userInfo
	user, ok := ctx.Value(userInfoKey).(interface{})
	if !ok {
		return nil
	}
	b, err := json.Marshal(user)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = json.Unmarshal(b, &usr)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return usr
}
