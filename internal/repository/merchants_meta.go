package repository

import (
	"cheemshappy_pay/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type MerchantsMetaRepository interface {
	GetMerchantsMeta(ctx context.Context, id int64) (*model.MerchantsMetum, error)
	CreateMerchantsMeta(ctx context.Context, req *model.Merchant) (*model.MerchantsMetum, error)
	UpdateMerchantsMeta(ctx context.Context, id int64, req *model.MerchantsMetum) error
	DeleteMerchantsMetaByMID(ctx context.Context, mid int64) error
}

func NewMerchantsMetaRepository(
	repository *Repository,
) MerchantsMetaRepository {
	return &merchantsMetaRepository{
		Repository: repository,
	}
}

type merchantsMetaRepository struct {
	*Repository
}

func (r *merchantsMetaRepository) GetMerchantsMeta(ctx context.Context, id int64) (*model.MerchantsMetum, error) {
	tx := newMerchantsMetum(r.DB(ctx))
	merchantsMeta, err := tx.Where(tx.MID.Eq(id)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return merchantsMeta, nil
}

func (r *merchantsMetaRepository) CreateMerchantsMeta(ctx context.Context, req *model.Merchant) (*model.MerchantsMetum, error) {
	merchantsMeta := &model.MerchantsMetum{
		MID: req.ID,
	}
	tx := newMerchantsMetum(r.DB(ctx))
	err := tx.Create(merchantsMeta)
	return merchantsMeta, err
}

func (r *merchantsMetaRepository) UpdateMerchantsMeta(ctx context.Context, id int64, req *model.MerchantsMetum) error {
	tx := newMerchantsMetum(r.DB(ctx))
	_, err := tx.Where(tx.MID.Eq(id)).Updates(req)
	return err
}

func (r *merchantsMetaRepository) DeleteMerchantsMetaByMID(ctx context.Context, mid int64) error {
	tx := newMerchantsMetum(r.DB(ctx))
	_, err := tx.Where(tx.MID.Eq(mid)).Delete()
	return err
}
