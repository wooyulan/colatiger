//go:build wireinject
// +build wireinject

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

var serverSet = wire.NewSet(server.NewHttpServer)

// build App
func newApp(httpServer *http.Server) *server.App {
	return server.NewApp(
		server.WithServer(httpServer),
		server.WithName("cola-tiger-server"),
	)
}

func NewWire(*config.Configuration, *zap.Logger, *lumberjack.Logger) (*server.App, func(), error) {
	panic(wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		middleware.ProviderSet,
		common.ProviderSet,
		serverSet,
		newApp,
	))
}
