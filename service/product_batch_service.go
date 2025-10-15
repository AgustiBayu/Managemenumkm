package service

import (
	"Managemenumkm/domain"
)

type ProductBatchService interface {
	FindByBarcode(barcode string) (domain.ProductBatchResponse, error)
	FindByProductID(productID uint) ([]domain.ProductBatchResponse, error)
	FindById(batchID uint) (domain.ProductBatchResponse, error)
	Update(request domain.ProductBatchUpdateRequest) (domain.ProductBatchResponse, error)
}

// Implementation is in product_batch_service_impl.go
