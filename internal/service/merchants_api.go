package service

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
)

type MerchantsApiService interface {
	GetMerchantsApi(ctx context.Context, id int64) (*model.MerchantsAPI, error)
	CreateMerchantsApi(ctx context.Context, req *v1.MerApiParam) error
	UpdateMerchantsApi(ctx context.Context, req *v1.MerApiParam) error
	DeleteMerchantsApi(ctx context.Context, id int64) error
	TestCallback(ctx context.Context, id int64) error
	GetMerchantsApiByApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error)
	CheckApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error)
	GetMerchantsApiList(ctx context.Context, req *v1.MerApiListReq) (*v1.Paginator, error)
	GetMerchantsApiByMIDs(ctx context.Context, mid int64) (*model.MerchantsAPI, error)
}

func NewMerchantsApiService(
	service *Service,
	merchantsApiRepository repository.MerchantsApiRepository,
) MerchantsApiService {
	return &merchantsApiService{
		Service:                service,
		merchantsApiRepository: merchantsApiRepository,
	}
}

type merchantsApiService struct {
	*Service
	merchantsApiRepository repository.MerchantsApiRepository
}

func (s *merchantsApiService) GetMerchantsApiList(ctx context.Context, req *v1.MerApiListReq) (*v1.Paginator, error) {
	return s.merchantsApiRepository.GetMerchantsApiList(ctx, req)
}

func (s *merchantsApiService) GetMerchantsApi(ctx context.Context, id int64) (*model.MerchantsAPI, error) {
	return s.merchantsApiRepository.GetMerchantsApi(ctx, id)
}

func (s *merchantsApiService) CreateMerchantsApi(ctx context.Context, req *v1.MerApiParam) error {
	//TODO: 生成apikey
	apikey := uuid.New().String()
	isExist, err := s.merchantsApiRepository.GetMerchantsApiByMID(ctx, req.MID, apikey)
	if err != nil {
		return err
	}
	if isExist != nil {
		return errors.New("apikey已存在,请重新生成")
	}

	ma := &model.MerchantsAPI{
		MID:         req.MID,
		Apikey:      apikey,
		CallbackURL: req.CallbackURL,
		SecretKey:   req.SecretKey,
		Remark:      req.Remark,
	}
	return s.merchantsApiRepository.CreateMerchantsApi(ctx, ma)
}

func (s *merchantsApiService) UpdateMerchantsApi(ctx context.Context, req *v1.MerApiParam) error {
	ma, err := s.merchantsApiRepository.GetMerchantsApi(ctx, req.ID)
	if err != nil {
		return err
	}
	ma.CallbackURL = req.CallbackURL
	ma.SecretKey = req.SecretKey
	ma.Remark = req.Remark
	return s.merchantsApiRepository.UpdateMerchantsApi(ctx, ma)
}

func (s *merchantsApiService) DeleteMerchantsApi(ctx context.Context, id int64) error {
	return s.merchantsApiRepository.DeleteMerchantsApi(ctx, id)
}

func (s *merchantsApiService) TestCallback(ctx context.Context, id int64) error {
	ma, err := s.merchantsApiRepository.GetMerchantsApi(ctx, id)
	if err != nil {
		return err
	}
	client := req.C()
	resp, err := client.R().
		SetHeader("c9pays", ma.SecretKey).
		SetBodyJsonMarshal(map[string]interface{}{
			"order_no": "1234567890",
			"amount":   100,
			"currency": "CNY",
			"status":   "success",
		}).
		Post(ma.CallbackURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("callback failed")
	}
	if resp.String() != "success" {
		return errors.New("callback failed not success")
	}
	return nil
}

func (s *merchantsApiService) GetMerchantsApiByApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error) {
	return s.merchantsApiRepository.GetMerchantsApiByApikey(ctx, apikey)
}

func (s *merchantsApiService) CheckApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error) {
	return s.merchantsApiRepository.CheckApikey(ctx, apikey)
}

func (s *merchantsApiService) GetMerchantsApiByMIDs(ctx context.Context, mid int64) (*model.MerchantsAPI, error) {
	return s.merchantsApiRepository.GetMerchantsApiByMIDs(ctx, mid)
}
