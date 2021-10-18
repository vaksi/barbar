package transport

import (
	"barbar/config"
	"barbar/domain/auth/repository"
	"barbar/domain/auth/usecase"
	authorization "barbar/pkg/middleware"
	"barbar/pkg/mongodb"
	"barbar/pkg/redis"
	"barbar/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

func RunHttp() {
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg)

	timeoutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// init mongo
	mongo, err := mongodb.Connect(timeoutCtx, cfg.Mongo.URL, cfg.Mongo.AuthDatabase)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Client().Disconnect(timeoutCtx)

	// init redis
	redisConn, err := redis.NewRedis("auth", cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// set dependency
	authRepo := repository.NewAuthRepository(mongo, redisConn)

	authUseCase := usecase.NewAuthUseCase(authRepo)

	authhandler := authHttpHandler{
		authUseCase: authUseCase,
	}

	authMiddleware := authorization.NewAuthMiddleware(redisConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger, utils.RenderJsonMiddleware, authMiddleware.CheckToken)
	r.Route(cfg.AuthHTTP.PathPrefix, func(r chi.Router) {
		r.Post("/login", authhandler.Login)
		r.Get("/logout", authhandler.Logout)
	})

	fmt.Println("Running User Service on PORT " + cfg.AuthHTTP.Port)
	http.ListenAndServe(cfg.AuthHTTP.Port, r)
}

type authHttpHandler struct {
	authUseCase usecase.AuthUseCaseInterface
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *authHttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.ErrorResponse(w, r, err)
		return
	}

	auth, err := l.authUseCase.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	render.JSON(w, r, auth)
	return
}

func (l *authHttpHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userInfo := authorization.GetUserInfoFromContext(r.Context())

	err := l.authUseCase.Logout(r.Context(), userInfo.UserID)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, fmt.Sprintf("user %s success logout", userInfo.Email))
	return
}
