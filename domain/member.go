package domain

import "time"

// Member represents a customer who is registered in the loyalty program
type Member struct {
	ID                uint       `gorm:"primaryKey"`
	TokoID            uint       `gorm:"not null"`
	Toko              Toko       `gorm:"foreignKey:TokoID"`
	MemberCode        string     `gorm:"uniqueIndex;not null;size:20"` // Unique member code
	Name              string     `gorm:"not null;size:100"`
	Email             string     `gorm:"size:100"`
	Phone             string     `gorm:"size:20"`
	Address           string     `gorm:"size:255"`
	Birthday          string     `gorm:"size:10"` // Format: YYYY-MM-DD
	Gender            string     `gorm:"size:10"` // Male, Female, Other
	Notes             string     `gorm:"size:500"`
	TotalPoints       int        `gorm:"default:0"`
	TotalSpent        float64    `gorm:"default:0"`              // Total amount spent historically
	TotalTransactions int        `gorm:"default:0"`              // Total number of transactions
	Status            string     `gorm:"default:active;size:20"` // active, inactive, suspended
	IsActive          bool       `gorm:"default:true"`           // Legacy field for compatibility
	JoinedDate        time.Time  `gorm:"not null"`
	LastVisitDate     *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// MemberTier defines different levels of membership benefits
type MemberTier struct {
	ID           uint    `gorm:"primaryKey"`
	TokoID       uint    `gorm:"not null"`
	Toko         Toko    `gorm:"foreignKey:TokoID"`
	Name         string  `gorm:"not null;size:50"` // Bronze, Silver, Gold, Platinum
	Description  string  `gorm:"size:255"`
	MinPoints    int     `gorm:"default:0"`     // Minimum points to reach this tier
	MinSpent     float64 `gorm:"default:0"`     // Minimum amount spent to reach this tier
	DiscountRate float64 `gorm:"default:0"`     // Discount rate (e.g., 0.05 for 5%)
	PointRate    float64 `gorm:"default:1"`     // Points earned per rupiah (e.g., 1 point per Rp 100)
	IsDefault    bool    `gorm:"default:false"` // Default tier for new members
	IsActive     bool    `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MemberTransaction records point transactions for members
type MemberTransaction struct {
	ID             uint         `gorm:"primaryKey"`
	TokoID         uint         `gorm:"not null"`
	Toko           Toko         `gorm:"foreignKey:TokoID"`
	MemberID       uint         `gorm:"not null"`
	Member         Member       `gorm:"foreignKey:MemberID"`
	TransactionID  *uint        // Link to sales transaction if applicable
	Transaction    *Transaction `gorm:"foreignKey:TransactionID"`
	Type           string       `gorm:"not null;size:20"` // EARN, REDEEM, EXPIRE, ADJUST
	Points         int          `gorm:"not null"`         // Positive for earning, negative for redemption
	Amount         float64      `gorm:"default:0"`        // Transaction amount in rupiah
	Description    string       `gorm:"size:255"`
	BalanceBefore  int          `gorm:"not null"` // Point balance before transaction
	BalanceAfter   int          `gorm:"not null"` // Point balance after transaction
	ExpirationDate *time.Time   // When points will expire
	IsExpired      bool         `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// PointRedemptionRule defines how points can be redeemed
type PointRedemptionRule struct {
	ID             uint     `gorm:"primaryKey"`
	TokoID         uint     `gorm:"not null"`
	Toko           Toko     `gorm:"foreignKey:TokoID"`
	Name           string   `gorm:"not null;size:100"`
	Type           string   `gorm:"not null;size:20"` // DISCOUNT, PRODUCT
	RequiredPoints int      `gorm:"not null"`
	DiscountRate   float64  `gorm:"default:0"` // For DISCOUNT type (e.g., 1000 points = 10% discount)
	DiscountAmount float64  `gorm:"default:0"` // For DISCOUNT type (e.g., 500 points = Rp 5000 off)
	ProductID      *uint    // For PRODUCT type
	Product        *Product `gorm:"foreignKey:ProductID"`
	MinPurchase    float64  `gorm:"default:0"` // Minimum purchase required
	MaxDiscount    float64  `gorm:"default:0"` // Maximum discount amount
	IsActive       bool     `gorm:"default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MemberFilter defines filter criteria for member queries
type MemberFilter struct {
	TokoID uint
	Search string // Search by name, phone, or member code
	Tier   string // Filter by tier name
	Status string // Filter by status (active, inactive, etc.)
}
