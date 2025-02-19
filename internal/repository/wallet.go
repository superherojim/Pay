package repository

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletByMID(ctx context.Context, mid int64) (*model.Wallet, error)
	GetWallet(ctx context.Context, req *v1.WalletList) (*v1.Paginator, error)
	DeleteWalletByMID(ctx context.Context, mid int64) error
	DeleteWallet(ctx context.Context, id int64) error
	CreateWallet(ctx context.Context, w *model.Wallet) error
	UpdateWallet(ctx context.Context, w *model.Wallet) error
}

func NewWalletRepository(
	repository *Repository,
) WalletRepository {
	return &walletRepository{
		Repository: repository,
	}
}

type walletRepository struct {
	*Repository
}

func (r *walletRepository) GetWalletByMID(ctx context.Context, mid int64) (*model.Wallet, error) {
	tx := newWallet(r.DB(ctx))
	wa, err := tx.Where(tx.MID.Eq(mid)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return wa, nil
}

func (r *walletRepository) GetWallet(ctx context.Context, req *v1.WalletList) (*v1.Paginator, error) {
	tx := newWallet(r.DB(ctx))
	if req.MID > 0 {
		tx.Where(tx.MID.Eq(req.MID))
	}
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

func (r *walletRepository) DeleteWalletByMID(ctx context.Context, mid int64) error {
	tx := newWallet(r.DB(ctx))
	_, err := tx.Where(tx.MID.Eq(mid)).Delete()
	return err
}

func (r *walletRepository) DeleteWallet(ctx context.Context, id int64) error {
	tx := newWallet(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(id)).Delete()
	return err
}

func (r *walletRepository) CreateWallet(ctx context.Context, w *model.Wallet) error {
	tx := newWallet(r.DB(ctx))
	wa, err := tx.Where(tx.MID.Eq(w.MID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if wa != nil && wa.ID > 0 {
		return errors.New("wallet already exists")
	}
	err = tx.Create(w)
	return err
}

func (r *walletRepository) UpdateWallet(ctx context.Context, w *model.Wallet) error {
	tx := newWallet(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(w.ID)).Updates(w)
	return err
}
