package main

import (
	"colatiger/cmd/wire"
	"colatiger/pkg/config"
	"colatiger/pkg/log"
	"context"
	"flag"
	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "../config/local.yaml", "config path, eg: -conf ./config/local.yaml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", "http://127.0.0.1:"+conf.GetString("http.port")))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
