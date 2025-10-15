package service

import (
	"Managemenumkm/domain"
)

// ProductService handles the business logic for master products.
type ProductService interface {
	Create(req domain.ProductCreateRequest) (domain.ProductResponse, error)
	Update(req domain.ProductUpdateRequest) (domain.ProductResponse, error)
	Delete(productID uint) error
	FindById(productID uint) (domain.ProductResponse, error)
	FindAll() ([]domain.ProductResponse, error)

	// AddStock creates a new batch for an existing product.
	AddStock(req domain.AddStockStep2Request) (domain.ProductBatchResponse, error)
}