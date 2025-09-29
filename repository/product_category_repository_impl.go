package repository

import (
	"Managemenumkm/domain"
	"context"

	"gorm.io/gorm"
)

type ProductCategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) ProductCategoryRepository {
	return &ProductCategoryRepositoryImpl{
		DB: db,
	}
}

func (p *ProductCategoryRepositoryImpl) Create(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
	if err := p.DB.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}
func (p *ProductCategoryRepositoryImpl) FindAll(ctx context.Context) ([]*domain.ProductCategory, error) {
	var product []*domain.ProductCategory
	if err := p.DB.WithContext(ctx).Find(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
func (p *ProductCategoryRepositoryImpl) FindById(ctx context.Context, categoryId int) (*domain.ProductCategory, error) {
	var category *domain.ProductCategory
	if err := p.DB.WithContext(ctx).First(&category, categoryId).Error; err != nil {
		return nil, err
	}
	return category, nil
}
func (p *ProductCategoryRepositoryImpl) Update(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error) {
	if err := p.DB.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}
func (p *ProductCategoryRepositoryImpl) Delete(ctx context.Context, category *domain.ProductCategory) error {
	if err := p.DB.WithContext(ctx).Delete(&category).Error; err != nil {
		return err
	}
	return nil
}
