package repository

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MerchantsApiRepository interface {
	GetMerchantsApi(ctx context.Context, id int64) (*model.MerchantsAPI, error)
	GetMerchantsApiByMID(ctx context.Context, mid int64, apikey string) (*model.MerchantsAPI, error)
	CreateMerchantsApi(ctx context.Context, merchantsApi *model.MerchantsAPI) error
	UpdateMerchantsApi(ctx context.Context, merchantsApi *model.MerchantsAPI) error
	DeleteMerchantsApi(ctx context.Context, id int64) error
	GetMerchantsApiByApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error)
	CheckApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error)
	GetMerchantsApiList(ctx context.Context, req *v1.MerApiListReq) (*v1.Paginator, error)
	GetMerchantsApiByMIDs(ctx context.Context, mid int64) (*model.MerchantsAPI, error)
}

func NewMerchantsApiRepository(
	repository *Repository,
	rdb *redis.Client,
) MerchantsApiRepository {
	return &merchantsApiRepository{
		Repository: repository,
		rdb:        rdb,
	}
}

type merchantsApiRepository struct {
	*Repository
	rdb *redis.Client
}

func (r *merchantsApiRepository) GetMerchantsApiList(ctx context.Context, req *v1.MerApiListReq) (*v1.Paginator, error) {
	tx := newMerchantsAPI(r.DB(ctx))
	if req.MID != 0 {
		tx.Where(tx.MID.Eq(req.MID))
	}
	byPage, i, err := tx.FindByPage(req.GetOffset(), req.Size)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p := &v1.Paginator{
		Total: i,
		Data:  byPage,
	}
	return p, nil
}

func (r *merchantsApiRepository) GetMerchantsApi(ctx context.Context, id int64) (*model.MerchantsAPI, error) {
	tx := newMerchantsAPI(r.DB(ctx))
	merchantsApi, err := tx.Where(tx.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return merchantsApi, nil
}

func (r *merchantsApiRepository) GetMerchantsApiByMID(ctx context.Context, mid int64, apikey string) (*model.MerchantsAPI, error) {
	tx := newMerchantsAPI(r.DB(ctx))
	merchantsApi, err := tx.Where(tx.MID.Eq(mid), tx.Apikey.Eq(apikey)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return merchantsApi, nil
}

func (r *merchantsApiRepository) CreateMerchantsApi(ctx context.Context, merchantsApi *model.MerchantsAPI) error {
	tx := newMerchantsAPI(r.DB(ctx))
	return tx.Create(merchantsApi)
}

func (r *merchantsApiRepository) UpdateMerchantsApi(ctx context.Context, merchantsApi *model.MerchantsAPI) error {
	tx := newMerchantsAPI(r.DB(ctx))
	return tx.Save(merchantsApi)
}

func (r *merchantsApiRepository) DeleteMerchantsApi(ctx context.Context, id int64) error {
	tx := newMerchantsAPI(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(id)).Delete()
	return err
}

func (r *merchantsApiRepository) GetMerchantsApiByApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error) {
	tx := newMerchantsAPI(r.DB(ctx))
	merchantsApi, err := tx.Where(tx.Apikey.Eq(apikey)).First()
	if err != nil {
		return nil, err
	}
	return merchantsApi, nil
}

func (r *merchantsApiRepository) CheckApikey(ctx context.Context, apikey string) (*model.MerchantsAPI, error) {
	ma, err := r.CheckApikeyCache(ctx, apikey)
	if err == nil {
		return ma, nil
	}
	ma, err = r.GetMerchantsApiByApikey(ctx, apikey)
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(ma)
	if err != nil {
		return nil, err
	}
	r.rdb.Set(ctx, apikey, string(json), time.Hour*24)
	return ma, nil
}

func (r *merchantsApiRepository) CheckApikeyCache(ctx context.Context, apikey string) (*model.MerchantsAPI, error) {

	s, err := r.rdb.Get(ctx, apikey).Result()
	if err != nil {
		return nil, err
	}
	if s == "" {
		return nil, errors.New("apikey not found")
	}
	var ma model.MerchantsAPI
	err = json.Unmarshal([]byte(s), &ma)
	if err != nil {
		return nil, err
	}
	return &ma, nil
}

func (r *merchantsApiRepository) GetMerchantsApiByMIDs(ctx context.Context, mid int64) (*model.MerchantsAPI, error) {
	tx := newMerchantsAPI(r.DB(ctx))
	tx.Where(tx.MID.Eq(mid))
	result, err := tx.First()
	if err != nil {
		return nil, err
	}
	return result, nil
}
