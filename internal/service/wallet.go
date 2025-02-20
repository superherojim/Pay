package service

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"context"
)

type WalletService interface {
	GetWallet(ctx context.Context, req *v1.WalletList) (*v1.Paginator, error)
	DeleteWallet(ctx context.Context, id int64) error
	DeleteWalletByMID(ctx context.Context, mid int64) error
	CreateWallet(ctx context.Context, req *v1.Wallet) error
	UpdateWallet(ctx context.Context, w *v1.AddWallet) error
	GetWalletByMID(ctx context.Context, mid int64) (*model.Wallet, error)
	AddWallet(ctx context.Context, w *v1.AddWallet) error
	CreateChildWallet(ctx context.Context, wallet *model.Wallet) error
}

func NewWalletService(
	repo repository.WalletRepository,
	sysRepo repository.SysWalletRepository,
) WalletService {
	return &walletService{
		repo:    repo,
		sysRepo: sysRepo,
	}
}

type walletService struct {
	repo    repository.WalletRepository
	sysRepo repository.SysWalletRepository
}

func (s *walletService) GetWallet(ctx context.Context, req *v1.WalletList) (*v1.Paginator, error) {
	return s.repo.GetWallet(ctx, req)
}

func (s *walletService) DeleteWallet(ctx context.Context, id int64) error {
	return s.repo.DeleteWallet(ctx, id)
}

func (s *walletService) DeleteWalletByMID(ctx context.Context, mid int64) error {
	return s.repo.DeleteWalletByMID(ctx, mid)
}

func (s *walletService) CreateWallet(ctx context.Context, req *v1.Wallet) error {
	return s.repo.CreateWallet(ctx, nil)
}

func (s *walletService) AddWallet(ctx context.Context, w *v1.AddWallet) error {
	wallet := &model.Wallet{
		MID:      w.MID,
		Ac:       w.Ac,
		PriKey:   w.PriKey,
		Mnemonic: w.Mnemonic,
		Path:     w.Path,
		Remark:   w.Remark,
	}
	return s.repo.CreateWallet(ctx, wallet)
}

func (s *walletService) UpdateWallet(ctx context.Context, w *v1.AddWallet) error {
	wallet := &model.Wallet{
		ID:       w.ID,
		MID:      w.MID,
		Ac:       w.Ac,
		PriKey:   w.PriKey,
		Mnemonic: w.Mnemonic,
		Path:     w.Path,
		Remark:   w.Remark,
	}
	return s.repo.UpdateWallet(ctx, wallet)
}

func (s *walletService) GetWalletByMID(ctx context.Context, mid int64) (*model.Wallet, error) {
	return s.repo.GetWalletByMID(ctx, mid)
}

func (s *walletService) CreateChildWallet(ctx context.Context, wallet *model.Wallet) error {
	return s.repo.CreateWallet(ctx, wallet)
}
