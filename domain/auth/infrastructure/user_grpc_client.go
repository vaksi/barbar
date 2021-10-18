package infrastructure

import (
	"barbar/config"
	"barbar/proto/users"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	logger "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

func NewUserService() users.UsersServiceClient {
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg)

	cc := createConn(cfg.GrpcClient.GRPCUserURL,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithUnaryInterceptor(unaryInterceptorUser),
	)

	return users.NewUsersServiceClient(cc)
}

// createConn : Function to create grpc dial connection
func createConn(endpoint string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	cc, err := grpc.Dial(endpoint, dialOptions...)
	if err != nil {
		logger.Error("Error when creating", zap.Error(err))
	} else {
		logger.Info("Successfully connect to GRPC Server: " + endpoint)
	}
	return cc
}

func unaryInterceptorUser(ctx context.Context,
	method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("interceptorUser_%s", method))
	defer span.Finish()

	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)
	if err != nil {
		logger.Error("User Interceptor Log",
			zap.Any("Request", req),
			zap.String("Method", method),
			zap.Error(err),
			zap.Duration("Elapsed", elapsed),
		)
	} else {
		logger.Info("User Interceptor Log",
			zap.Any("Request", req),
			zap.String("Method", method),
			zap.Any("Response", reply),
			zap.Duration("Elapsed", elapsed),
		)
	}
	return err
}
