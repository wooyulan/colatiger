package server

import (
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
	userHandler handler.UserHandler) *http.Server {

	gin.SetMode(gin.ReleaseMode)

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	s.Use(
		middleware.CORSMiddleware(),
		//middleware.ResponseLogMiddleware(logger),
		//middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	//s.GET("/", func(ctx *gin.Context) {
	//	logger.WithContext(ctx).Info("hello")
	//	apiV1.HandleSuccess(ctx, map[string]interface{}{
	//		":)": "Thank you for using nunu!",
	//	})
	//})

	v1 := s.Group("/api/v1")
	noAuthRouter := v1
	{
		noAuthRouter.POST("/register", userHandler.Register)
		//noAuthRouter.POST("/login", userHandler.Login)
	}

	// Non-strict permission routing group
	//noStrictAuthRouter := v1.Use(middleware.NoStrictAuth(jwt, logger))
	//{
	//	noStrictAuthRouter.GET("/user", userHandler.GetProfile)
	//}
	//
	//// Strict permission routing group
	//strictAuthRouter := v1.Use(middleware.StrictAuth(jwt, logger))
	//{
	//	strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
	//}

	return s

}
