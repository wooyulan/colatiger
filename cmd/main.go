package main

import (
	"colatiger/cmd/bootstrap"
	"colatiger/cmd/wire"
	"context"
	"flag"
	"go.uber.org/zap"
	"strconv"
)

func main() {
	var envConf = flag.String("conf", "./conf/local.yaml", "conf path, eg: -conf ./conf/local.yaml")
	flag.Parse()

	config := bootstrap.NewConfig(*envConf)
	logger, loggerWriter := bootstrap.NewLog(config)

	app, cleanup, err := wire.NewWire(config, logger, loggerWriter)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", config.App.AppUrl+":"+strconv.Itoa(config.App.Port)))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
