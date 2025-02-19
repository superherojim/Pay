package service

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"bk/internal/repository"
	"context"
	"errors"
)

type MerchantsService interface {
	GetMerchants(ctx context.Context, id int64) (*model.Merchant, error)
	GetMerchantsList(ctx context.Context, req *v1.MerchantsListReq) (*v1.Paginator, error)
	CreateMerchants(ctx context.Context, req *v1.MerParam) error
	DeleteMerchants(ctx context.Context, id int64) error
	UpdateMerchants(ctx context.Context, req *v1.MerParam) error
	GetMerchantsIN(ctx context.Context) ([]*v1.MerIN, error)
}

func NewMerchantsService(
	service *Service,
	merchantsRepository repository.MerchantsRepository,
	merchantsMetaService MerchantsMetaService,
	walletService WalletService,
) MerchantsService {
	return &merchantsService{
		Service:              service,
		merchantsRepository:  merchantsRepository,
		merchantsMetaService: merchantsMetaService,
		walletService:        walletService,
	}
}

type merchantsService struct {
	*Service
	merchantsRepository  repository.MerchantsRepository
	merchantsMetaService MerchantsMetaService
	walletService        WalletService
}

func (s *merchantsService) GetMerchants(ctx context.Context, id int64) (*model.Merchant, error) {
	return s.merchantsRepository.GetMerchants(ctx, id)
}
func (s *merchantsService) GetMerchantsList(ctx context.Context, req *v1.MerchantsListReq) (*v1.Paginator, error) {
	merchants, err := s.merchantsRepository.GetMerchantsList(ctx, req)
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

func (s *merchantsService) CreateMerchants(ctx context.Context, req *v1.MerParam) error {
	merchants, err := s.merchantsRepository.GetMerchantsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if merchants != nil {
		return errors.New("用户已存在")
	}
	mer, err := s.merchantsRepository.CreateMerchants(ctx, req)
	if err != nil {
		return err
	}
	_, err = s.merchantsMetaService.CreateMerchantsMeta(ctx, mer)
	if err != nil {
		return err
	}
	return nil
}

func (s *merchantsService) DeleteMerchants(ctx context.Context, id int64) error {
	err := s.merchantsRepository.DeleteMerchants(ctx, id)
	if err != nil {
		return err
	}
	err = s.merchantsMetaService.DeleteMerchantsMetaByMID(ctx, id)
	if err != nil {
		return err
	}
	err = s.walletService.DeleteWalletByMID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *merchantsService) UpdateMerchants(ctx context.Context, req *v1.MerParam) error {
	mer, err := s.merchantsRepository.GetMerchantsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if mer != nil && mer.ID != req.ID {
		return errors.New("用户已存在")
	}
	if req.Avatar != "" {
		mer.Avatar = req.Avatar
	}
	if req.Nickname != "" {
		mer.Nickname = req.Nickname
	}
	if req.Introduction != "" {
		mer.Introduction = req.Introduction
	}
	if req.Phone != "" {
		mer.Phone = req.Phone
	}
	err = s.merchantsRepository.UpdateMerchants(ctx, mer)
	if err != nil {
		return err
	}
	return nil
}

func (s *merchantsService) GetMerchantsIN(ctx context.Context) ([]*v1.MerIN, error) {
	merchants, err := s.merchantsRepository.GetMerchantsIN(ctx)
	if err != nil {
		return nil, err
	}
	return merchants, nil
}
