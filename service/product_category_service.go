package service

import (
	"Managemenumkm/domain"
	"context"
)

type ProductCategoryService interface {
	Create(ctx context.Context, req *domain.ProductCategoryCreateRequest) error
	FindAll(ctx context.Context) ([]*domain.ProductCategoryResponse, error)
	FindById(ctx context.Context, categoryID int) (*domain.ProductCategoryResponse, error)
	Update(ctx context.Context, req *domain.ProductCategoryUpdateRequest) error
	Delete(ctx context.Context, categoryID int) error
}
