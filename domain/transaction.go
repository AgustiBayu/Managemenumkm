package domain

import "time"

// Transaction represents the header of a sales transaction
type Transaction struct {
	ID              uint              `gorm:"primaryKey"`
	TokoID          uint              `gorm:"not null"`
	Toko            Toko              `gorm:"foreignKey:TokoID"`
	UserID          uint              `gorm:"not null"`
	User            User              `gorm:"foreignKey:UserID"`
	TransactionDate time.Time         `gorm:"not null"`
	TotalAmount     float64           `gorm:"not null"`
	PaymentMethod   string            `gorm:"not null;size:50"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TransactionItem represents a single product line within a transaction
type TransactionItem struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID uint      `gorm:"not null"`
	ProductID     uint      `gorm:"not null"`
	Product       Product   `gorm:"foreignKey:ProductID"`
	Quantity      int       `gorm:"not null"`
	Price         float64   `gorm:"not null"` // Price at the time of transaction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
