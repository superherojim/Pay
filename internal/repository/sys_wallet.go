package repository

import (
	"bk/internal/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SysWalletRepository interface {
	Get(ctx context.Context) (*model.SysWallet, error)
	Create(ctx context.Context, wallet *model.SysWallet) error
	Update(ctx context.Context, wallet *model.SysWallet) error
	UpdateIndex(ctx context.Context, index uint32) error
	IncrementIndex(ctx context.Context) (uint32, error)
	Exists(ctx context.Context) (bool, error)
}

func NewSysWalletRepository(repository *Repository) SysWalletRepository {
	return &sysWalletRepository{
		Repository: repository,
	}
}

type sysWalletRepository struct {
	*Repository
}

func (r *sysWalletRepository) Get(ctx context.Context) (*model.SysWallet, error) {
	tx := newSysWallet(r.DB(ctx))
	wallet, err := tx.First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return wallet, nil
}

func (r *sysWalletRepository) Create(ctx context.Context, wallet *model.SysWallet) error {
	tx := newSysWallet(r.DB(ctx))
	return tx.Create(wallet)
}

func (r *sysWalletRepository) Update(ctx context.Context, wallet *model.SysWallet) error {
	tx := newSysWallet(r.DB(ctx))
	_, err := tx.Where(tx.ID.Eq(wallet.ID)).Updates(wallet)
	return err
}

func (r *sysWalletRepository) UpdateIndex(ctx context.Context, index uint32) error {
	return r.db.WithContext(ctx).
		Model(&model.SysWallet{}).
		Update("current_index", index).
		Error
}

func (r *sysWalletRepository) IncrementIndex(ctx context.Context) (uint32, error) {
	var currentIndex int32
	var version int32
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 使用SELECT FOR UPDATE锁定行并获取版本号
		var sysWallet model.SysWallet
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("current_index, version").
			Where("id > ?", 0).
			First(&sysWallet).
			Error; err != nil {
			return err
		}
		version = sysWallet.Version
		currentIndex = sysWallet.CurrentIndex
		newIndex := currentIndex + 1
		// 添加版本号检查防止并发更新
		result := tx.Model(&model.SysWallet{}).
			Where("version = ?", version).
			Updates(map[string]interface{}{
				"current_index": newIndex,
				"version":       version + 1,
			})

		if result.RowsAffected == 0 {
			return fmt.Errorf("并发修改冲突，请重试")
		}
		return result.Error
	})
	return uint32(currentIndex + 1), err
}

// 新增检查系统钱包存在的方法
func (r *sysWalletRepository) checkSysWalletExists(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.SysWallet{}).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sysWalletRepository) Exists(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.SysWallet{}).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
