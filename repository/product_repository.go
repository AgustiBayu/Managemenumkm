package repository

import (
	"Managemenumkm/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Save(product domain.Product) (domain.Product, error)
	FindAll() ([]domain.Product, error)
	FindById(productID uint) (domain.Product, error)
	FindBySKU(sku string) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(productID uint) error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{DB: db}
}