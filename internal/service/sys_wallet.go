package service

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"cheemshappy_pay/pkg/wallet"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type SysWalletService interface {
	GetSysWallet(ctx context.Context) (*model.SysWallet, error)
	UpdateSysWallet(ctx context.Context, req *v1.SysWalletUpdateReq) error
	CreateSysWallet(ctx context.Context) error
	GetMasterWallet(ctx context.Context) (*wallet.MasterWalletInfo, error)
	DeriveChildWallet(ctx context.Context) (*model.Wallet, error)
}

func NewSysWalletService(
	repo repository.SysWalletRepository,
) SysWalletService {
	return &sysWalletService{
		repo: repo,
	}
}

type sysWalletService struct {
	repo repository.SysWalletRepository
}

func (s *sysWalletService) GetSysWallet(ctx context.Context) (*model.SysWallet, error) {
	return s.repo.Get(ctx)
}

func (s *sysWalletService) UpdateSysWallet(ctx context.Context, req *v1.SysWalletUpdateReq) error {
	wallet := &model.SysWallet{
		Ac:       req.Ac,
		PriKey:   req.PriKey,
		Mnemonic: req.Mnemonic,
		Path:     req.Path,
		Remark:   req.Remark,
	}
	return s.repo.Update(ctx, wallet)
}

func (s *sysWalletService) CreateSysWallet(ctx context.Context) error {
	// 创建主钱包
	masterInfo, err := wallet.CreateSystemMasterWallet()
	if err != nil {
		return err
	}

	// 保存到数据库
	sysWallet := &model.SysWallet{
		Ac:       masterInfo.Address,
		PriKey:   masterInfo.PrivateKey,
		Mnemonic: masterInfo.Mnemonic,
		Path:     masterInfo.DerivePath,
	}
	return s.repo.Create(ctx, sysWallet)
}

func (s *sysWalletService) GetMasterWallet(ctx context.Context) (*wallet.MasterWalletInfo, error) {
	sysWallet, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &wallet.MasterWalletInfo{
		Mnemonic:   sysWallet.Mnemonic,
		Address:    sysWallet.Ac,
		PrivateKey: sysWallet.PriKey,
		DerivePath: sysWallet.Path,
	}, nil
}

func (s *sysWalletService) DeriveChildWallet(ctx context.Context) (*model.Wallet, error) {
	// 确保系统钱包存在
	exists, err := s.repo.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("检查系统钱包失败: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("请先创建系统钱包")
	}

	// 获取当前索引并自增
	index, err := s.repo.IncrementIndex(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to increment index: %v", err)
	}

	// 获取主钱包信息
	masterInfo, err := s.GetMasterWallet(ctx)
	if err != nil {
		return nil, err
	}

	// 派生子钱包
	childWallet, err := wallet.DeriveChildWallet(masterInfo, index)
	if err != nil {
		return nil, err
	}

	return &model.Wallet{
		Ac:     childWallet.Address.Hex(),
		PriKey: hex.EncodeToString(crypto.FromECDSA(childWallet.PrivateKey)),
		Path:   childWallet.Path,
	}, nil
}
