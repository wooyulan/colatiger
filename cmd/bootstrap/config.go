package bootstrap

import (
	"colatiger/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func NewConfig(path string) *config.Configuration {
	return initConf(path)
}

func initConf(configPath string) *config.Configuration {
	fmt.Println("load conf:" + configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read conf failed: %s \n", err))
	}

	var conf *config.Configuration

	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("conf file changed:", in.Name)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		if err := v.Unmarshal(&conf); err != nil {
			fmt.Println(err)
		}
	})
	return conf
}
