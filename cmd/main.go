package main

import (
	"colatiger/cmd/wire"
	"colatiger/pkg/conf"
	"colatiger/pkg/log"
	"context"
	"flag"
	"go.uber.org/zap"
	"strconv"
)

func main() {
	var envConf = flag.String("conf", "./conf/local.yaml", "conf path, eg: -conf ./conf/local.yaml")
	flag.Parse()

	config := conf.NewConfig(*envConf)
	logger := log.NewLog(config)

	app, cleanup, err := wire.NewWire(config, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", config.App.AppUrl+strconv.Itoa(config.App.Port)))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
