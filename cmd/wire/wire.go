//go:build wireinject
// +build wireinject

package wire

import (
	"colatiger/internal/handler"
	"colatiger/internal/repository"
	"colatiger/internal/server"
	"colatiger/internal/service"
	"colatiger/pkg/app"
	"colatiger/pkg/helper/sid"
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"colatiger/pkg/server/http"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewChatHandler,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)
var serverSet = wire.NewSet(
	server.NewHttpServer,
)

// build App
func newApp(httpServer *http.Server) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("cola-tiger-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		jwt.NewJwt,
		sid.NewSid,
		newApp,
	))
}
