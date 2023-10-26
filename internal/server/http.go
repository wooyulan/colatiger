package server

import (
	v1 "colatiger/api/v1/res"
	"colatiger/internal/handler"
	"colatiger/internal/middleware"
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"colatiger/pkg/server/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewHttpServer(logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler handler.UserHandler,
	chatHandler handler.ChatHandler,
) *http.Server {

	// 初始化验证器
	middleware.InitializeValidator()

	// 初始化表结构
	middleware.InitializeDB(conf)

	gin.SetMode(gin.ReleaseMode)

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	//s.Use(
	//	middleware.CORSMiddleware(),
	//)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		v1.Success(ctx, "welcome user colatiger")
	})

	s.POST("/chat", chatHandler.ChatStream)

	v1 := s.Group("/api/v1")
	noAuthRouter := v1
	{
		noAuthRouter.POST("/auth/register", userHandler.Register)
		noAuthRouter.POST("/auth/login", userHandler.Login)
	}

	// Non-strict permission routing group
	authRouter := v1.Use(middleware.NoStrictAuth(jwt))
	{
		authRouter.GET("/user/info", userHandler.GetInfo)
	}
	return s

}
