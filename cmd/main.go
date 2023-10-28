package main

import (
	"colatiger/cmd/wire"
	"colatiger/pkg/conf"
	"colatiger/pkg/log"
	"context"
	"flag"
	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "./conf/local.yaml", "conf path, eg: -conf ./conf/local.yaml")
	flag.Parse()

	configGlobal := conf.NewConfig(*envConf)
	logger := log.NewLog(configGlobal)

	app, cleanup, err := wire.NewWire(configGlobal, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", "http://127.0.0.1:"+conf.GetString("http.port")))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
