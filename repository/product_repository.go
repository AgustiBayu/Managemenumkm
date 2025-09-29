package repository

import (
	"Managemenumkm/domain"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	FindAll(ctx context.Context) ([]*domain.Product, map[uint]*domain.ProductCategory, error)
	FindById(ctx context.Context, produkId uint) (*domain.Product, *domain.ProductCategory, error)
	FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	FindLowStock(ctx context.Context, threshold uint) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) (*domain.Product, error)
	UpdateStock(ctx context.Context, productId uint, newStock uint) error
	Delete(ctx context.Context, product *domain.Product) error
	UploadThumbnail(ctx context.Context, productId uint, path string) error
}
