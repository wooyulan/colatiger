package log

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logger struct {
	*logrus.Logger
}

func NewLog(conf *viper.Viper) *Logger {
	return InitializeLog(conf)
}
func InitializeLog(conf *viper.Viper) *Logger {
	return &Logger{logrus.New()}
}
