// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"colatiger/config"
	"colatiger/internal/handler"
	"colatiger/internal/middleware"
	"colatiger/internal/repository"
	"colatiger/internal/server"
	"colatiger/internal/service"
	"colatiger/pkg/common"
	"colatiger/pkg/server/http"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Injectors from wire.go:

func NewWire(configuration *config.Configuration, logger *zap.Logger, lumberjackLogger *lumberjack.Logger) (*server.App, func(), error) {
	cors := middleware.NewCors()
	db := repository.NewDB(configuration, logger)
	client := repository.NewRedis(configuration, logger)
	sonyflake := common.NewSonyFlake()
	minioClient := repository.NewOss(configuration, logger)
	clientClient := repository.NewMilvus(configuration, logger)
	repositoryRepository, cleanup, err := repository.NewRepository(logger, db, client, sonyflake, minioClient, clientClient)
	if err != nil {
		return nil, nil, err
	}
	userRepo := repository.NewUserRepository(logger, repositoryRepository)
	userService := service.NewUserService(userRepo)
	lockBuilder := common.NewLockBuilder(client)
	jwtRepo := repository.NewJwtRepo(repositoryRepository, logger)
	jwtService := service.NewJwtService(configuration, logger, userService, lockBuilder, jwtRepo)
	jwtAuth := middleware.NewJWTAuth(configuration, jwtService)
	authHandler := handler.NewAuthHandler(logger, jwtService, userService)
	milvusRepo := repository.NewMilvusRepository(repositoryRepository, logger)
	milvusService := service.NewMilvusService(milvusRepo)
	chatRepo := repository.NewChatRepository(logger, repositoryRepository)
	chatService := service.NewChatService(chatRepo)
	chatHandler := handler.NewChatHandler(logger, milvusService, chatService)
	ossHandler := handler.NewOssHandler(logger, minioClient, sonyflake, configuration)
	recovery := middleware.NewRecovery(lumberjackLogger)
	ocrRepo := repository.NewOcrRepository(logger, repositoryRepository)
	ocrService := service.NewOcrService(ocrRepo)
	ocrHandler := handler.NewOcrHandler(ocrService)
	httpServer := server.NewHttpServer(logger, configuration, cors, jwtAuth, authHandler, chatHandler, ossHandler, recovery, ocrHandler)
	app := newApp(httpServer)
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

var serverSet = wire.NewSet(server.NewHttpServer)

// build App
func newApp(httpServer *http.Server) *server.App {
	return server.NewApp(server.WithServer(httpServer), server.WithName("cola-tiger-server"))
}
