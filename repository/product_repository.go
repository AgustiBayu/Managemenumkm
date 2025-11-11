package repository

import (
	"Managemenumkm/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Save(product domain.Product) (domain.Product, error)
	FindAll() ([]domain.Product, error)
	FindById(productID uint, tx ...*gorm.DB) (domain.Product, error)
	FindBySKU(sku string) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(productID uint) error
	UpdateBatchStock(tx *gorm.DB, batchID uint, newStock uint) error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{DB: db}
}