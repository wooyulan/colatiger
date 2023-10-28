package conf

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
