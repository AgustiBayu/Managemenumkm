package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product mendefinisikan informasi umum dan master dari sebuah produk.
type Product struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	SKU        string `gorm:"unique"`
	Price      uint
	ImageURL   string // URL to the product image
	CategoryID uint
	DeletedAt  gorm.DeletedAt  `gorm:"index"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID"`
	Batches    []ProductBatch  `gorm:"foreignKey:ProductID"` // Relasi ke batch
}

// ProductBatch merepresentasikan setiap batch/lot stok yang masuk.
type ProductBatch struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	Barcode   string `gorm:"unique"`
	Exp       time.Time
	Stock     uint
	DateAdded time.Time
}

// --- Request & Response Structs ---

type ProductCreateRequest struct {
	Name       string `validate:"required" form:"name"`
	SKU        string `validate:"required" form:"sku"`
	Price      uint   `validate:"required" form:"price"`
	CategoryID uint   `validate:"required" form:"category_id"`
}

type ProductUpdateRequest struct {
	ID         uint   `validate:"required" form:"id"`
	Name       string `validate:"required" form:"name"`
	SKU        string `validate:"required" form:"sku"`
	Price      uint   `validate:"required" form:"price"`
	CategoryID uint   `validate:"required" form:"category_id"`
}

type ProductResponse struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	SKU             string                  `json:"sku"`
	Price           uint                    `json:"price"`
	ImageURL        string                  `json:"image_url"`
	TotalStock      uint                    `json:"total_stock"` // Agregasi dari semua batch
	CategoryID      uint                    `json:"category_id"`
	ProductCategory ProductCategoryResponse `json:"category"`
}

type ProductBatchResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Barcode   string    `json:"barcode"`
	Exp       string    `json:"exp"`
	Stock     uint      `json:"stock"`
	DateAdded string    `json:"date_added"`
}

// AddStockStep1Request untuk mencari barcode
type AddStockStep1Request struct {
	Barcode string `form:"barcode" validate:"required"`
}

// AddStockStep2Request untuk membuat batch baru jika barcode tidak ditemukan
type AddStockStep2Request struct {
	Barcode   string `form:"barcode" validate:"required"`
	ProductID uint   `form:"product_id" validate:"required"`
	Exp       string `form:"exp" validate:"required"`
	Stock     uint   `form:"stock" validate:"required,gt=0"`
}

// ProductBatchUpdateRequest for updating an existing batch
type ProductBatchUpdateRequest struct {
	ID        uint   `validate:"required"`
	ProductID uint   `validate:"required"`
	Stock     uint   `validate:"required,min=0"`
	Exp       string `validate:"required"`
}