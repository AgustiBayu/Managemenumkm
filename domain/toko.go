package domain

import (
	"time"

	"gorm.io/gorm"
)

type Toko struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Address   string
	Users     []User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TokoCreateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type TokoUpdateRequest struct {
	ID      uint   `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type TokoResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
