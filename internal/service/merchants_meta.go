package service

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"context"
)

type MerchantsMetaService interface {
	GetMerchantsMeta(ctx context.Context, id int64) (*model.MerchantsMetum, error)
	CreateMerchantsMeta(ctx context.Context, req *model.Merchant) (*model.MerchantsMetum, error)
	UpdateMerchantsMeta(ctx context.Context, id int64, req *v1.MerchantsMeta) error
	DeleteMerchantsMetaByMID(ctx context.Context, mid int64) error
}

func NewMerchantsMetaService(
	service *Service,
	merchantsMetaRepository repository.MerchantsMetaRepository,
) MerchantsMetaService {
	return &merchantsMetaService{
		Service:                 service,
		merchantsMetaRepository: merchantsMetaRepository,
	}
}

type merchantsMetaService struct {
	*Service
	merchantsMetaRepository repository.MerchantsMetaRepository
}

func (s *merchantsMetaService) GetMerchantsMeta(ctx context.Context, id int64) (*model.MerchantsMetum, error) {
	return s.merchantsMetaRepository.GetMerchantsMeta(ctx, id)
}

func (s *merchantsMetaService) CreateMerchantsMeta(ctx context.Context, req *model.Merchant) (*model.MerchantsMetum, error) {
	return s.merchantsMetaRepository.CreateMerchantsMeta(ctx, req)
}

func (s *merchantsMetaService) UpdateMerchantsMeta(ctx context.Context, id int64, req *v1.MerchantsMeta) error {
	mm := &model.MerchantsMetum{
		Ac: req.Ac,
	}
	return s.merchantsMetaRepository.UpdateMerchantsMeta(ctx, id, mm)
}

func (s *merchantsMetaService) DeleteMerchantsMetaByMID(ctx context.Context, mid int64) error {
	return s.merchantsMetaRepository.DeleteMerchantsMetaByMID(ctx, mid)
}
