package transport

import (
	"barbar/config"
	"barbar/domain/users/repository"
	"barbar/domain/users/usecase"
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
	mongo, err := mongodb.Connect(timeoutCtx, cfg.Mongo.URL, cfg.Mongo.UserDatabase)
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
	userRepo := repository.NewUserRepository(mongo, redisConn)

	userUseCase := usecase.NewUserUseCase(userRepo)

	userHandler := userHttpHandler{
		userUseCase: userUseCase,
	}

	authorizationMiddleware := authorization.NewAuthMiddleware(redisConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger, utils.RenderJsonMiddleware, authorizationMiddleware.CheckToken)
	r.Route(cfg.UserHTTP.PathPrefix, func(r chi.Router) {
		r.Post("/users", userHandler.RegisterUser)
		r.Put("/users/{uid}", userHandler.UpdateUser)
		r.Delete("/users/{uid}", userHandler.DeleteUser)
		r.Get("/users/{uid}", userHandler.GetUserById)
		r.Get("/users", userHandler.GetAllUser)
	})

	fmt.Println("Running User Service on PORT " + cfg.UserHTTP.Port)
	http.ListenAndServe(cfg.UserHTTP.Port, r)
}

type userHttpHandler struct {
	userUseCase usecase.UserUseCaseInterface
}

type userRequest struct {
	Name     *string `json:"name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (h *userHttpHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request userRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	user, err := h.userUseCase.Register(r.Context(), *request.Name, *request.Phone, *request.Email, *request.Password)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
	return
}

func (h *userHttpHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var request userRequest

	uid := chi.URLParam(r, "uid")

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	user, err := h.userUseCase.Update(r.Context(), uid, request.Name, request.Phone, request.Email, request.Password)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
	return
}

func (h *userHttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	err := h.userUseCase.Delete(r.Context(), uid)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	render.JSON(w, r, "user deleted")
	return
}

func (h *userHttpHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	user, err := h.userUseCase.GetById(r.Context(), uid)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, user)
	return
}

func (h *userHttpHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.GetAllUser(r.Context(), "", map[string]bool{
		"_id": false,
	}, 10, 0)
	if err != nil {
		logger.Error(err)
		utils.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, users)
	return
}