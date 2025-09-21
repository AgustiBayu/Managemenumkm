package repository

import (
	"Managemenumkm/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TokoRepositoryImpl struct {
	DB *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &TokoRepositoryImpl{
		DB: db,
	}
}

func (r *TokoRepositoryImpl) Create(ctx context.Context, toko *domain.Toko) (*domain.Toko, error) {
	if err := r.DB.WithContext(ctx).Create(toko).Error; err != nil {
		return nil, fmt.Errorf("create toko failed err: %s", err)
	}
	return toko, nil
}

func (r *TokoRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*domain.Toko, int64, error) {
	var toko []*domain.Toko
	var totalItem int64

	if err := r.DB.WithContext(ctx).Model(&domain.Toko{}).Count(&totalItem).Error; err != nil {
		return nil, 0, fmt.Errorf("toko is null: %s", err)
	}
	offset := (page - 1) * pageSize
	if err := r.DB.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&toko).Error; err != nil {
		return nil, 0, err
	}
	return toko, totalItem, nil
}

func (r *TokoRepositoryImpl) FindById(ctx context.Context, tokoID int) (*domain.Toko, error) {
	var toko *domain.Toko
	if err := r.DB.WithContext(ctx).First(&toko, tokoID).Error; err != nil {
		return nil, err
	}
	return toko, nil
}

func (r *TokoRepositoryImpl) Update(ctx context.Context, toko *domain.Toko) (*domain.Toko, error) {
	if err := r.DB.WithContext(ctx).Save(&toko).Error; err != nil {
		return nil, fmt.Errorf("failed update toko err: %s", err)
	}
	return toko, nil
}

func (r *TokoRepositoryImpl) Delete(ctx context.Context, toko *domain.Toko) error {
	if err := r.DB.WithContext(ctx).Delete(&toko).Error; err != nil {
		return fmt.Errorf("failed delete toko err: %s", err)
	}
	return nil
}
