package repository

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type MerchantsRepository interface {
	GetMerchants(ctx context.Context, id int64) (*model.Merchant, error)
	CreateMerchants(ctx context.Context, req *v1.MerParam) (*model.Merchant, error)
	GetMerchantsByEmail(ctx context.Context, email string) (*model.Merchant, error)
	DeleteMerchants(ctx context.Context, id int64) error
	UpdateMerchants(ctx context.Context, model *model.Merchant) error
	GetMerchantsList(ctx context.Context, req *v1.MerchantsListReq) (*v1.Paginator, error)
	GetMerchantsIN(ctx context.Context) ([]*v1.MerIN, error)
	TotalCount(ctx context.Context) (int64, error)
}

func NewMerchantsRepository(
	repository *Repository,
) MerchantsRepository {
	return &merchantsRepository{
		Repository: repository,
	}
}

type merchantsRepository struct {
	*Repository
}

func (r *merchantsRepository) GetMerchants(ctx context.Context, id int64) (*model.Merchant, error) {
	tx := newMerchant(r.DB(ctx))
	merchants, err := tx.Where(tx.ID.Eq(id)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return merchants, nil
}

func (r *merchantsRepository) GetMerchantsList(ctx context.Context, req *v1.MerchantsListReq) (*v1.Paginator, error) {
	tx := newMerchant(r.DB(ctx))
	tx.Order(tx.CreatedAt.Desc())
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

func (r *merchantsRepository) CreateMerchants(ctx context.Context, req *v1.MerParam) (*model.Merchant, error) {
	merchants := &model.Merchant{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}
	tx := newMerchant(r.DB(ctx))
	err := tx.Create(merchants)
	return merchants, err
}

func (r *merchantsRepository) GetMerchantsByEmail(ctx context.Context, email string) (*model.Merchant, error) {
	tx := newMerchant(r.DB(ctx))
	u, err := tx.Where(tx.Email.Eq(email)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return u, nil
}

func (r *merchantsRepository) DeleteMerchants(ctx context.Context, id int64) error {
	tx := newMerchant(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(id)).Delete()
	return err
}

func (r *merchantsRepository) UpdateMerchants(ctx context.Context, model *model.Merchant) error {
	tx := newMerchant(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(model.ID)).Updates(model)
	return err
}

func (r *merchantsRepository) GetMerchantsIN(ctx context.Context) ([]*v1.MerIN, error) {
	tx := newMerchant(r.DB(ctx))
	tx.Select(tx.ID, tx.Nickname)
	tx.Order(tx.CreatedAt.Desc())
	m := make([]*v1.MerIN, 0)
	ls, err := tx.Find()
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		m = append(m, &v1.MerIN{
			ID:       v.ID,
			Nickname: v.Nickname,
		})
	}
	return m, nil
}

func (r *merchantsRepository) TotalCount(ctx context.Context) (int64, error) {
	tx := newMerchant(r.DB(ctx))
	return tx.Count()
}
