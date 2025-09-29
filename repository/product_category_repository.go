package repository

import (
	"Managemenumkm/domain"
	"context"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error)
	FindAll(ctx context.Context) ([]*domain.ProductCategory, error)
	FindById(ctx context.Context, categoryId int) (*domain.ProductCategory, error)
	Update(ctx context.Context, category *domain.ProductCategory) (*domain.ProductCategory, error)
	Delete(ctx context.Context, category *domain.ProductCategory) error
}
