package service

import (
	"cheemshappy_pay/internal/repository"
	"cheemshappy_pay/pkg/helper/sid"
	"cheemshappy_pay/pkg/jwt"
	"cheemshappy_pay/pkg/log"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(tm repository.Transaction, logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
