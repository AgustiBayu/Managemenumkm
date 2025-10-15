package repository

import (
	"Managemenumkm/domain"

	"gorm.io/gorm"
)

type ProductBatchRepository interface {
	Save(batch domain.ProductBatch) (domain.ProductBatch, error)
	FindByBarcode(barcode string) (domain.ProductBatch, error)
	FindByProductID(productID uint) ([]domain.ProductBatch, error)
	FindById(batchID uint) (domain.ProductBatch, error)
	Update(batch domain.ProductBatch) (domain.ProductBatch, error)
}

func NewProductBatchRepository(db *gorm.DB) ProductBatchRepository {
	return &productBatchRepositoryImpl{DB: db}
}
