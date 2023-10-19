package handler

import (
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
)

type Handler struct {
	logger *log.Logger
	jwt    *jwt.JWT
}

func NewHandler(logger *log.Logger, jwt *jwt.JWT) *Handler {
	return &Handler{
		logger: logger,
		jwt:    jwt,
	}
}
