package middleware

import (
	"colatiger/api/v1/res"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Recovery struct {
	loggerWriter *lumberjack.Logger
}

func NewRecoveryM(loggerWriter *lumberjack.Logger) *Recovery {
	return &Recovery{
		loggerWriter: loggerWriter,
	}
}

func (m *Recovery) Handler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		m.loggerWriter,
		res.ServerError,
	)
}
