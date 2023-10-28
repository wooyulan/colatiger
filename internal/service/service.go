package service

import (
	"github.com/google/wire"
)

//type Service struct {
//	logger *log.Logger
//	sid    *sid.Sid
//}
//
//func NewService(logger *log.Logger, sid *sid.Sid) *Service {
//	return &Service{
//		logger: logger,
//		sid:    sid,
//	}
//}

var Providerset = wire.NewSet(NewUserService)
