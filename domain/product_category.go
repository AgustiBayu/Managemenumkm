package domain

import "gorm.io/gorm"

type ProductCategory struct {
	ID        uint `gorm:"primaryKey"`
	Category  string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProductCategoryCreateRequest struct {
	Category string `validate:"required" json:"category"`
}

type ProductCategoryResponse struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
}

type ProductCategoryUpdateRequest struct {
	ID       uint   `validate:"required" json:"id"`
	Category string `validate:"required" json:"category"`
}
