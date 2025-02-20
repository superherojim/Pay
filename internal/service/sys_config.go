package service

import (
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"context"
)

type SysConfigService interface {
	GetSysConfig(ctx context.Context) (*model.SysConfig, error)
	CreateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error
	UpdateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error
}

func NewSysConfigService(
	service *Service,
	sysConfigRepository repository.SysConfigRepository,
) SysConfigService {
	return &sysConfigService{
		Service:             service,
		sysConfigRepository: sysConfigRepository,
	}
}

type sysConfigService struct {
	*Service
	sysConfigRepository repository.SysConfigRepository
}

func (s *sysConfigService) GetSysConfig(ctx context.Context) (*model.SysConfig, error) {
	return s.sysConfigRepository.GetSysConfig(ctx)
}

func (s *sysConfigService) CreateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error {
	return s.sysConfigRepository.CreateSysConfig(ctx, sysConfig)
}

func (s *sysConfigService) UpdateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error {
	return s.sysConfigRepository.UpdateSysConfig(ctx, sysConfig)
}
