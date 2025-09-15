package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 uint `gorm:"primaryKey"`
	FName              string
	LName              string
	Email              string `gorm:"uniq"`
	Password           string
	Alamat             string
	Thumbnail          string
	Role               string
	SubscriptionStatus string
	SubscriptionExpiry time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type UserUpdateRequest struct {
	ID                 uint      `json:"id"`
	FName              string    `json:"first_name" validate:"requaired"`
	LName              string    `json:"last_name"`
	Email              string    `json:"email" validate:"requaired"`
	Password           string    `json:"password" validate:"requaired"`
	Alamat             string    `json:"alamat" validate:"requaired"`
	Thumbnail          string    `json:"thumbnail"`
	Role               string    `json:"role" validate:"requaired"`
	SubscriptionStatus string    `json:"subscription_status" validate:"requaired"`
	SubscriptionExpiry time.Time `json:"subscription_expriry" validate:"requaired"`
	CreatedAt          time.Time `json:"created_ad" validate:"requaired"`
	UpdatedAt          time.Time `json:"updated_ad" validate:"requaired"`
}

type UserCreateRequest struct {
	FName              string    `json:"first_name" validate:"requaired"`
	LName              string    `json:"last_name"`
	Email              string    `json:"email" validate:"requaired"`
	Password           string    `json:"password" validate:"requaired"`
	Alamat             string    `json:"alamat" validate:"requaired"`
	Thumbnail          string    `json:"thumbnail"`
	Role               string    `json:"role" validate:"requaired"`
	SubscriptionStatus string    `json:"subscription_status" validate:"requaired"`
	SubscriptionExpiry time.Time `json:"subscription_expiry" validate:"requaired"`
	CreatedAt          time.Time `json:"created_at" validate:"requaired"`
	UpdatedAt          time.Time `json:"updated_at" validate:"requaired"`
}

type UserResponse struct {
	ID                 uint      `json:"id"`
	FName              string    `json:"first_name" validate:"requaired"`
	LName              string    `json:"last_name"`
	Email              string    `json:"email" validate:"requaired"`
	Alamat             string    `json:"alamat" validate:"requaired"`
	Thumbnail          string    `json:"thumbnail"`
	Role               string    `json:"role" validate:"requaired"`
	SubscriptionStatus string    `json:"subscription_status" validate:"requaired"`
	SubscriptionExpiry time.Time `json:"subscription_expiry" validate:"requaired"`
	CreatedAt          time.Time `json:"created_at" validate:"requaired"`
	UpdatedAt          time.Time `json:"updated_at" validate:"requaired"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
