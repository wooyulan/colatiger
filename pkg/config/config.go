package config

import (
	"colatiger/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//func NewConfig(p string) *viper.Viper {
//	envConf := os.Getenv("APP_CONF")
//	if envConf == "" {
//		envConf = p
//	}
//	fmt.Println("load conf file:", envConf)
//	return getConfig(envConf)
//}
//
//func getConfig(path string) *viper.Viper {
//	conf := viper.New()
//	conf.SetConfigFile(path)
//	err := conf.ReadInConfig()
//	if err != nil {
//		panic(err)
//	}
//	return conf
//}

func InitConf(configPath string) *config.Configuration {

	fmt.Println("load config:" + configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	var conf *config.Configuration

	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
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
