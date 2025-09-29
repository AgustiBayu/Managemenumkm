package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Barcode    string `gorm:"unique"`
	Thumbnail  string
	Price      uint
	Exp        time.Time
	Stock      uint
	CategoryID uint
	DeletedAt  gorm.DeletedAt  `gorm:"index"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID"`
}

type ProductCreateRequest struct {
	Name       string `validate:"required" form:"name"`
	Barcode    string `validate:"required" form:"barcode"`
	Thumbnail  string
	Price      uint   `validate:"required" form:"price"`
	Exp        string `validate:"required" form:"exp"`
	Stock      uint   `validate:"required" form:"stock"`
	CategoryID uint   `validate:"required" form:"category_id"`
}

type ProductResponse struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	Barcode         string                  `json:"barcode"`
	Thumbnail       string                  `json:"thumbnail"`
	Price           uint                    `json:"price"`
	Exp             string                  `json:"exp"`
	Stock           uint                    `json:"stock"`
	CategoryID      uint                    `json:"category_id"`
	ProductCategory ProductCategoryResponse `json:"category"`
}

type ProductUpdateRequest struct {
	ID         uint   `validate:"required" form:"id"`
	Name       string `validate:"required" form:"name"`
	Barcode    string `validate:"required" form:"barcode"`
	Thumbnail  string
	Price      uint   `validate:"required" form:"price"`
	Exp        string `validate:"required" form:"exp"`
	Stock      uint   `validate:"required" form:"stock"`
	CategoryID uint   `validate:"required" form:"category_id"`
}

type ProductUpdateStockRequest struct {
	Stock uint `json:"stock" validate:"required,gte=0"`
}