package transport

import (
	"barbar/config"
	"barbar/domain/users/repository"
	handler "barbar/domain/users/transport/grpc"
	"barbar/domain/users/usecase"
	"barbar/pkg/mongodb"
	"barbar/pkg/redis"
	"barbar/proto/users"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
	logger "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func ServeGRPC() {
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg)
	grpcListener := initGRPCListener(cfg)
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

	userServer := handler.NewGrpcServer(userUseCase)
	//create grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	users.RegisterUsersServiceServer(grpcServer, userServer)
	if err := grpcServer.Serve(grpcListener); err != nil {
		logger.Error("Error Serve GRPC", zap.Error(err))
	}
}

func initGRPCListener(cfg *config.MainConfig) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", cfg.GrpcServer.UserPort))
	if err != nil {
		logger.Error("Net Listen Failed", zap.Error(err))
	}
	logger.Info(fmt.Sprintf("GRPC Service Online at: 0.0.0.0%s", cfg.GrpcServer.UserPort))
	return lis
}

func unaryInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	ext.SamplingPriority.Set(span, 1)
	defer span.Finish()
	start := time.Now()

	var errStatus string
	res, err := handler(ctx, req)

	elapsed := time.Since(start)
	if err != nil {
		span.SetTag(string(ext.Error), true)
		logger.WithFields(logger.Fields{
			"Request": req,
			"Error":   err,
			"Elapsed": elapsed,
		}).Errorln("Interceptor Log")
		errStatus = "error"
	} else {
		logger.Info(false, "grpc-user", info.FullMethod, "GRPC User Interceptor log", req, res)
	}
	span.LogFields(
		tlog.String("status", errStatus),
		tlog.String("operation", info.FullMethod),
		tlog.String("message", "Debt Notes Interceptor log"),
		tlog.String("request", fmt.Sprintf("%v", req)),
		tlog.String("response", fmt.Sprintf("%v", res)),
		tlog.String("error_message", fmt.Sprintf("%v", err)),
	)

	return res, err
}
