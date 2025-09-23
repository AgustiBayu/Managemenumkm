package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 uint `gorm:"primaryKey"`
	FName              string
	LName              string
	Email              string `gorm:"unique"`
	Password           string
	Alamat             string
	Thumbnail          string
	Number             string
	Role               Role `gorm:"default:cashier"`
	TokoID             uint
	Toko               Toko `gorm:"foreignKey:TokoID"`
	SubscriptionStatus string
	SubscriptionExpiry time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type UserUpdateRequest struct {
	ID                 uint   `json:"id"`
	FName              string `json:"first_name" validate:"required"`
	LName              string `json:"last_name" validate:"required"`
	Email              string `json:"email" validate:"required"`
	Password           string `json:"password"`
	Alamat             string `json:"alamat" validate:"required"`
	Thumbnail          string `json:"thumbnail"`
	Number             string `json:"number" validate:"required"`
	Role               Role   `json:"role" validate:"required"`
	TokoID             uint   `json:"toko_id" validate:"required"`
	SubscriptionStatus string `json:"subscription_status" validate:"required"`
	SubscriptionExpiry string `json:"subscription_expriry" validate:"required"`
}

type UserCreateRequest struct {
	FName              string `json:"first_name" validate:"required"`
	LName              string `json:"last_name" validate:"required"`
	Email              string `json:"email" validate:"required"`
	Password           string `json:"password" validate:"required"`
	Alamat             string `json:"alamat" validate:"required"`
	Thumbnail          string `json:"thumbnail"`
	Number             string `json:"number" validate:"required"`
	Role               Role   `json:"role" validate:"required"`
	TokoID             uint   `json:"toko_id" validate:"required"`
	SubscriptionStatus string `json:"subscription_status" validate:"required"`
	SubscriptionExpiry string `json:"subscription_expiry" validate:"required"`
}

type UserResponse struct {
	ID                     uint          `json:"id"`
	FName                  string        `json:"first_name"`
	LName                  string        `json:"last_name"`
	Email                  string        `json:"email"`
	Alamat                 string        `json:"alamat"`
	Thumbnail              string        `json:"thumbnail"`
	Number                 string        `json:"number"`
	Role                   string        `json:"role"`
	Toko                   *TokoResponse `json:"toko"`
	SubscriptionStatus     string        `json:"subscription_status"`
	SubscriptionExpiry     string        `json:"subscription_expiry"`
	CreatedAt              time.Time     `json:"created_at"`
	SubscriptionExpiryDate time.Time     `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
