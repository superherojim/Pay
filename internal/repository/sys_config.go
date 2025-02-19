package repository

import (
	"bk/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type SysConfigRepository interface {
	GetSysConfig(ctx context.Context) (*model.SysConfig, error)
	CreateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error
	UpdateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error
}

func NewSysConfigRepository(
	repository *Repository,
) SysConfigRepository {
	return &sysConfigRepository{
		Repository: repository,
	}
}

type sysConfigRepository struct {
	*Repository
}

func (r *sysConfigRepository) GetSysConfig(ctx context.Context) (*model.SysConfig, error) {
	tx := newSysConfig(r.db)
	sysConfig, err := tx.First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return sysConfig, nil
}

func (r *sysConfigRepository) CreateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error {
	return r.db.Create(sysConfig).Error
}

func (r *sysConfigRepository) UpdateSysConfig(ctx context.Context, sysConfig *model.SysConfig) error {
	return r.db.Save(sysConfig).Error
}
