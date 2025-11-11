package service

import (
	"Managemenumkm/domain"
	"context"
)

// ProductService handles the business logic for master products.
type ProductService interface {
	Create(request domain.ProductCreateRequest) (domain.ProductResponse, error)
	Update(request domain.ProductUpdateRequest) (domain.ProductResponse, error)
	Delete(productID uint) error
	FindById(productID uint) (domain.ProductResponse, error)
	FindAll() ([]domain.ProductResponse, error)
	AddStock(request domain.AddStockStep2Request) (domain.ProductBatchResponse, error)
	UpdateImageURL(ctx context.Context, productID uint, imageURL string) error
}
