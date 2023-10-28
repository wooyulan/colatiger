//go:build wireinject
// +build wireinject

package wire

import (
	"colatiger/internal/handler"
	"colatiger/internal/middleware"
	"colatiger/internal/repository"
	"colatiger/internal/server"
	"colatiger/internal/service"
	"colatiger/pkg/helper/sid"
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"colatiger/pkg/server/http"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var serverSet = wire.NewSet(server.NewHttpServer)

// build App
func newApp(httpServer *http.Server) *server.App {
	return server.NewApp(
		server.WithServer(httpServer),
		server.WithName("cola-tiger-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*server.App, func(), error) {
	panic(wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		middleware.ProviderSet,
		serverSet,
		jwt.NewJwt,
		sid.NewSid,
		newApp,
	))
}
