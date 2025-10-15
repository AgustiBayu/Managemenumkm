package repository

import (
	"Managemenumkm/domain"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	DB *gorm.DB
}

func (r *productRepositoryImpl) Save(product domain.Product) (domain.Product, error) {
	err := r.DB.Create(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepositoryImpl) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	// Eager load Category and Batches
	err := r.DB.Preload("Category").Preload("Batches").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepositoryImpl) FindById(productID uint) (domain.Product, error) {
	var product domain.Product
	// Eager load Category and Batches
	err := r.DB.Preload("Category").Preload("Batches").First(&product, productID).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepositoryImpl) FindBySKU(sku string) (domain.Product, error) {
	var product domain.Product
	err := r.DB.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepositoryImpl) Update(product domain.Product) (domain.Product, error) {
	err := r.DB.Save(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepositoryImpl) Delete(productID uint) error {
	return r.DB.Delete(&domain.Product{}, productID).Error
}