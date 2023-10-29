package middleware

import (
	"colatiger/api/response"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Recovery struct {
	loggerWriter *lumberjack.Logger
}

func NewRecovery(loggerWriter *lumberjack.Logger) *Recovery {
	return &Recovery{
		loggerWriter: loggerWriter,
	}
}

func (m *Recovery) Handler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		m.loggerWriter,
		response.ServerError,
	)
}
