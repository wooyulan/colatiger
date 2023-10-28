package server

import (
	v1 "colatiger/api/response"
	"colatiger/internal/handler"
	"colatiger/internal/middleware"
	"colatiger/pkg/log"
	"colatiger/pkg/server/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewHttpServer(logger *log.Logger,
	conf *viper.Viper,
	cors *middleware.Cors,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
) *http.Server {

	// 初始化验证器
	middleware.InitializeValidator()

	// 初始化表结构
	middleware.InitializeDB(conf)

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	if conf.GetString("app.env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

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
	}

	// Non-strict permission routing group
	authRouter := v1.Use()
	{
		authRouter.GET("/user/info", authHandler.GetInfo)
	}

	// 对话
	chatRouter := v1
	{
		chatRouter.POST("/chat/stream", middleware.HeadersMiddleware(), chatHandler.ChatStream)
	}

	return s

}
