package repository

import (
	"Managemenumkm/domain"

	"gorm.io/gorm"
)

type productBatchRepositoryImpl struct {
	DB *gorm.DB
}

func (r *productBatchRepositoryImpl) Save(batch domain.ProductBatch) (domain.ProductBatch, error) {
	err := r.DB.Create(&batch).Error
	if err != nil {
		return batch, err
	}
	return batch, nil
}

func (r *productBatchRepositoryImpl) FindByBarcode(barcode string) (domain.ProductBatch, error) {
	var batch domain.ProductBatch
	err := r.DB.Where("barcode = ?", barcode).First(&batch).Error
	if err != nil {
		return batch, err
	}
	return batch, nil
}

func (r *productBatchRepositoryImpl) FindByProductID(productID uint) ([]domain.ProductBatch, error) {
	var batches []domain.ProductBatch
	err := r.DB.Where("product_id = ?", productID).Find(&batches).Error
	if err != nil {
		return batches, err
	}
	return batches, nil
}

func (r *productBatchRepositoryImpl) FindById(batchID uint) (domain.ProductBatch, error) {
	var batch domain.ProductBatch
	err := r.DB.Where("id = ?", batchID).First(&batch).Error
	if err != nil {
		return batch, err
	}
	return batch, nil
}

func (r *productBatchRepositoryImpl) Update(batch domain.ProductBatch) (domain.ProductBatch, error) {
	err := r.DB.Save(&batch).Error
	if err != nil {
		return batch, err
	}
	return batch, nil
}
