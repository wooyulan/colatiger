//go:build wireinject
// +build wireinject

package wire

import (
	"colatiger/cmd"
	"colatiger/internal/server"
	"colatiger/pkg/helper/sid"
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"colatiger/pkg/server/http"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var serverSet = wire.NewSet(server.NewHttpServer)

// build App
func newApp(httpServer *http.Server) *main.App {
	return main.NewApp(
		main.WithServer(httpServer),
		main.WithName("cola-tiger-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*main.App, func(), error) {
	panic(wire.Build(
		//repository.ProviderSet,
		//service.Providerset,
		//handler.ProviderSet,
		//server.ProviderSet,
		serverSet,
		jwt.NewJwt,
		sid.NewSid,
		newApp,
	))
}
