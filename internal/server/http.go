package server

import (
	v1 "colatiger/api/response"
	"colatiger/config"
	"colatiger/internal/handler"
	"colatiger/internal/middleware"
	"colatiger/internal/model"
	"colatiger/pkg/server/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewHttpServer(logger *zap.Logger,
	conf *config.Configuration,
	cors *middleware.Cors,
	jwtAuth *middleware.JWTAuth,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	ossHandler *handler.OssHandler,
	recovery *middleware.Recovery,
) *http.Server {

	// 初始化验证器
	middleware.InitializeValidator()

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.App.AppUrl),
		http.WithServerPort(conf.App.Port),
	)

	if conf.App.Env == "prod" || conf.App.Env == "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.Use(gin.Logger(), recovery.Handler())
	//跨域处理
	s.Use(cors.CORSMiddleware())

	s.GET("/", func(ctx *gin.Context) {
		v1.Success(ctx, "welcome user colatiger")
	})

	v1 := s.Group("/api/v1")
	noAuthRouter := v1
	{
		noAuthRouter.POST("/auth/register", authHandler.Register)
		noAuthRouter.POST("/auth/login", authHandler.Login)
		noAuthRouter.POST("/upload", ossHandler.Upload)
	}

	// 对话
	chatRouter := v1
	{
		chatRouter.POST("/chat/stream", middleware.HeadersMiddleware(), chatHandler.ChatStream)
		chatRouter.GET("/milvus", chatHandler.Test)
	}

	// Non-strict permission routing group
	authRouter := v1.Use(jwtAuth.Handler(model.AppGuardName))
	{
		authRouter.GET("/user/info", authHandler.GetInfo)
	}

	return s

}
