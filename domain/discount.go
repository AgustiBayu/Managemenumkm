package domain

import "time"

// Discount represents various discount campaigns and promotions
type Discount struct {
	ID                uint                   `gorm:"primaryKey"`
	TokoID            uint                   `gorm:"not null"`
	Toko              Toko                   `gorm:"foreignKey:TokoID"`
	Name              string                 `gorm:"not null;size:100"`
	Description       string                 `gorm:"size:255"`
	Type              string                 `gorm:"not null;size:20"` // PERCENTAGE, FIXED_AMOUNT, BUY_X_GET_Y, BULK
	DiscountValue     float64                `gorm:"not null"`        // Discount amount or percentage
	MinPurchase       float64                `gorm:"default:0"`       // Minimum purchase amount
	MaxDiscount       float64                `gorm:"default:0"`       // Maximum discount amount (for percentage discounts)
	MinQuantity       int                    `gorm:"default:0"`       // Minimum quantity for bulk discounts
	BuyQuantity       int                    `gorm:"default:0"`       // Quantity to buy for BUY_X_GET_Y
	GetQuantity       int                    `gorm:"default:0"`       // Quantity to get for BUY_X_GET_Y
	GetDiscountPercent float64               `gorm:"default:0"`       // Discount percentage for free items in BUY_X_GET_Y
	ApplicableTo      string                 `gorm:"default:ALL;size:20"` // ALL, CATEGORY, PRODUCT, MEMBER_TIER
	ProductID         *uint                  // For specific product discounts
	Product           *Product               `gorm:"foreignKey:ProductID"`
	CategoryID        *uint                  // For category discounts
	Category          *ProductCategory       `gorm:"foreignKey:CategoryID"`
	MemberTierID      *uint                  // For member tier specific discounts
	MemberTier        *MemberTier            `gorm:"foreignKey:MemberTierID"`
	StartDate         time.Time              `gorm:"not null"`
	EndDate           time.Time              `gorm:"not null"`
	UsageLimit        int                    `gorm:"default:0"`       // 0 = unlimited
	UsageCount        int                    `gorm:"default:0"`       // How many times it has been used
	IsActive          bool                   `gorm:"default:true"`
	IsStackable       bool                   `gorm:"default:false"`   // Can be combined with other discounts
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// DiscountUsage tracks how discounts are used in transactions
type DiscountUsage struct {
	ID           uint         `gorm:"primaryKey"`
	TokoID       uint         `gorm:"not null"`
	Toko         Toko         `gorm:"foreignKey:TokoID"`
	DiscountID   uint         `gorm:"not null"`
	Discount     Discount     `gorm:"foreignKey:DiscountID"`
	TransactionID uint        `gorm:"not null"`
	Transaction  Transaction  `gorm:"foreignKey:TransactionID"`
	OriginalAmount float64    `gorm:"not null"` // Original item/subtotal amount
	DiscountAmount float64    `gorm:"not null"` // Amount discounted
	FinalAmount   float64     `gorm:"not null"` // Final amount after discount
	CreatedAt     time.Time
}

// SpecialOffer represents time-limited special promotions
type SpecialOffer struct {
	ID          uint      `gorm:"primaryKey"`
	TokoID      uint      `gorm:"not null"`
	Toko        Toko      `gorm:"foreignKey:TokoID"`
	Title       string    `gorm:"not null;size:100"`
	Description string    `gorm:"size:500"`
	BannerURL   string    `gorm:"size:255"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	Priority    int       `gorm:"default:0"` // Higher number = higher priority
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SpecialOfferProduct links products to special offers
type SpecialOfferProduct struct {
	ID              uint          `gorm:"primaryKey"`
	SpecialOfferID  uint          `gorm:"not null"`
	SpecialOffer    SpecialOffer  `gorm:"foreignKey:SpecialOfferID"`
	ProductID       uint          `gorm:"not null"`
	Product         Product       `gorm:"foreignKey:ProductID"`
	DiscountType    string        `gorm:"not null;size:20"` // PERCENTAGE, FIXED_AMOUNT
	DiscountValue   float64       `gorm:"not null"`
	MinQuantity     int           `gorm:"default:0"`
	MaxDiscount     float64       `gorm:"default:0"`
	CreatedAt       time.Time
}

// FlashSale represents limited-time flash sale events
type FlashSale struct {
	ID              uint      `gorm:"primaryKey"`
	TokoID          uint      `gorm:"not null"`
	Toko            Toko      `gorm:"foreignKey:TokoID"`
	Name            string    `gorm:"not null;size:100"`
	Description     string    `gorm:"size:255"`
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	StockLimit      int       `gorm:"default:0"` // 0 = unlimited
	UserLimit       int       `gorm:"default:1"` // Max purchase per user
	IsActive        bool      `gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FlashSaleProduct links products to flash sales
type FlashSaleProduct struct {
	ID             uint       `gorm:"primaryKey"`
	FlashSaleID    uint       `gorm:"not null"`
	FlashSale      FlashSale  `gorm:"foreignKey:FlashSaleID"`
	ProductID      uint       `gorm:"not null"`
	Product        Product    `gorm:"foreignKey:ProductID"`
	OriginalPrice  float64    `gorm:"not null"`
	SalePrice      float64    `gorm:"not null"`
	StockAvailable int        `gorm:"default:0"` // Available stock for this flash sale
	StockSold      int        `gorm:"default:0"`
	CreatedAt      time.Time
}

// --- Request & Response Structs ---

type DiscountCreateRequest struct {
	Name              string    `json:"name" validate:"required"`
	Description       string    `json:"description"`
	Type              string    `json:"type" validate:"required,oneof=PERCENTAGE FIXED_AMOUNT BUY_X_GET_Y BULK"`
	DiscountValue     float64   `json:"discount_value" validate:"required,gt=0"`
	MinPurchase       float64   `json:"min_purchase"`
	MaxDiscount       float64   `json:"max_discount"`
	MinQuantity       int       `json:"min_quantity"`
	BuyQuantity       int       `json:"buy_quantity"`
	GetQuantity       int       `json:"get_quantity"`
	GetDiscountPercent float64  `json:"get_discount_percent"`
	ApplicableTo      string    `json:"applicable_to" validate:"required,oneof=ALL CATEGORY PRODUCT MEMBER_TIER"`
	ProductID         *uint     `json:"product_id"`
	CategoryID        *uint     `json:"category_id"`
	MemberTierID      *uint     `json:"member_tier_id"`
	StartDate         string    `json:"start_date" validate:"required"`
	EndDate           string    `json:"end_date" validate:"required"`
	UsageLimit        int       `json:"usage_limit"`
	IsStackable       bool      `json:"is_stackable"`
}

type DiscountUpdateRequest struct {
	ID                uint      `json:"id" validate:"required"`
	Name              string    `json:"name" validate:"required"`
	Description       string    `json:"description"`
	Type              string    `json:"type" validate:"required,oneof=PERCENTAGE FIXED_AMOUNT BUY_X_GET_Y BULK"`
	DiscountValue     float64   `json:"discount_value" validate:"required,gt=0"`
	MinPurchase       float64   `json:"min_purchase"`
	MaxDiscount       float64   `json:"max_discount"`
	MinQuantity       int       `json:"min_quantity"`
	BuyQuantity       int       `json:"buy_quantity"`
	GetQuantity       int       `json:"get_quantity"`
	GetDiscountPercent float64  `json:"get_discount_percent"`
	ApplicableTo      string    `json:"applicable_to" validate:"required,oneof=ALL CATEGORY PRODUCT MEMBER_TIER"`
	ProductID         *uint     `json:"product_id"`
	CategoryID        *uint     `json:"category_id"`
	MemberTierID      *uint     `json:"member_tier_id"`
	StartDate         string    `json:"start_date" validate:"required"`
	EndDate           string    `json:"end_date" validate:"required"`
	UsageLimit        int       `json:"usage_limit"`
	IsStackable       bool      `json:"is_stackable"`
	IsActive          bool      `json:"is_active"`
}

type DiscountResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Type              string    `json:"type"`
	DiscountValue     float64   `json:"discount_value"`
	MinPurchase       float64   `json:"min_purchase"`
	MaxDiscount       float64   `json:"max_discount"`
	MinQuantity       int       `json:"min_quantity"`
	BuyQuantity       int       `json:"buy_quantity"`
	GetQuantity       int       `json:"get_quantity"`
	GetDiscountPercent float64  `json:"get_discount_percent"`
	ApplicableTo      string    `json:"applicable_to"`
	ProductID         *uint     `json:"product_id"`
	ProductName       *string   `json:"product_name"`
	CategoryID        *uint     `json:"category_id"`
	CategoryName      *string   `json:"category_name"`
	MemberTierID      *uint     `json:"member_tier_id"`
	MemberTierName    *string   `json:"member_tier_name"`
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date"`
	UsageLimit        int       `json:"usage_limit"`
	UsageCount        int       `json:"usage_count"`
	RemainingUsage    int       `json:"remaining_usage"`
	IsActive          bool      `json:"is_active"`
	IsStackable       bool      `json:"is_stackable"`
	CreatedAt         string    `json:"created_at"`
	UpdatedAt         string    `json:"updated_at"`
}

type SpecialOfferCreateRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	BannerURL   string    `json:"banner_url"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     string    `json:"end_date" validate:"required"`
	Priority    int       `json:"priority"`
	Products    []SpecialOfferProductRequest `json:"products" validate:"required,min=1"`
}

type SpecialOfferProductRequest struct {
	ProductID     uint    `json:"product_id" validate:"required"`
	DiscountType  string  `json:"discount_type" validate:"required,oneof=PERCENTAGE FIXED_AMOUNT"`
	DiscountValue float64 `json:"discount_value" validate:"required,gt=0"`
	MinQuantity   int     `json:"min_quantity"`
	MaxDiscount   float64 `json:"max_discount"`
}

type FlashSaleCreateRequest struct {
	Name        string                     `json:"name" validate:"required"`
	Description string                     `json:"description"`
	StartDate   string                     `json:"start_date" validate:"required"`
	EndDate     string                     `json:"end_date" validate:"required"`
	StockLimit  int                        `json:"stock_limit"`
	UserLimit   int                        `json:"user_limit"`
	Products    []FlashSaleProductRequest  `json:"products" validate:"required,min=1"`
}

type FlashSaleProductRequest struct {
	ProductID      uint    `json:"product_id" validate:"required"`
	SalePrice      float64 `json:"sale_price" validate:"required,gt=0"`
	StockAvailable int     `json:"stock_available"`
}

// --- Filter Structs ---

type DiscountFilter struct {
	TokoID       uint
	Type         string
	ApplicableTo string
	IsActive     *bool
	Search       string
}

type FlashSaleFilter struct {
	TokoID   uint
	IsActive *bool
	Search   string
}

// --- Helper Functions ---

func (d *Discount) IsExpired() bool {
	return time.Now().After(d.EndDate)
}

func (d *Discount) IsNotStarted() bool {
	return time.Now().Before(d.StartDate)
}

func (d *Discount) IsValid() bool {
	return d.IsActive && !d.IsExpired() && !d.IsNotStarted() && (d.UsageLimit == 0 || d.UsageCount < d.UsageLimit)
}

func (d *Discount) GetRemainingUsage() int {
	if d.UsageLimit == 0 {
		return -1 // Unlimited
	}
	return d.UsageLimit - d.UsageCount
}

func (fs *FlashSale) IsExpired() bool {
	return time.Now().After(fs.EndDate)
}

func (fs *FlashSale) IsNotStarted() bool {
	return time.Now().Before(fs.StartDate)
}

func (fs *FlashSale) IsValid() bool {
	return fs.IsActive && !fs.IsExpired() && !fs.IsNotStarted()
}

// BonusPointRule defines special conditions for earning bonus points
type BonusPointRule struct {
	ID              uint      `gorm:"primaryKey"`
	TokoID          uint      `gorm:"not null"`
	Toko            Toko      `gorm:"foreignKey:TokoID"`
	Name            string    `gorm:"not null;size:100"`
	Description     string    `gorm:"size:255"`
	Type            string    `gorm:"not null;size:20"` // MULTIPLIER, FIXED_BONUS, CATEGORY_BONUS, PRODUCT_BONUS, BIRTHDAY, FIRST_PURCHASE
	Multiplier      float64   `gorm:"default:1"`        // Points multiplier (e.g., 2x points = 2.0)
	BonusPoints     int       `gorm:"default:0"`         // Fixed bonus points
	MinPurchase     float64   `gorm:"default:0"`         // Minimum purchase amount
	ProductID       *uint     // For product-specific bonuses
	Product         *Product  `gorm:"foreignKey:ProductID"`
	CategoryID      *uint     // For category-specific bonuses
	Category        *ProductCategory `gorm:"foreignKey:CategoryID"`
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	UsageLimit      int       `gorm:"default:0"` // 0 = unlimited
	UsageCount      int       `gorm:"default:0"`
	IsActive        bool      `gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TierPromotion defines temporary requirements for member tier upgrades
type TierPromotion struct {
	ID                    uint       `gorm:"primaryKey"`
	TokoID                uint       `gorm:"not null"`
	Toko                  Toko       `gorm:"foreignKey:TokoID"`
	Name                  string     `gorm:"not null;size:100"`
	Description           string     `gorm:"size:255"`
	TargetTierID          uint       `gorm:"not null"`
	TargetTier            MemberTier `gorm:"foreignKey:TargetTierID"`
	PromotionalMinPoints  int        `gorm:"default:0"` // Temporary lower points requirement
	PromotionalMinSpent   float64    `gorm:"default:0"` // Temporary lower spending requirement
	OriginalMinPoints     int        `gorm:"default:0"` // Original points requirement
	OriginalMinSpent      float64    `gorm:"default:0"` // Original spending requirement
	StartDate             time.Time  `gorm:"not null"`
	EndDate               time.Time  `gorm:"not null"`
	MaxPromotions         int        `gorm:"default:0"` // Maximum members that can use this promotion
	UsedPromotions        int        `gorm:"default:0"`
	IsActive              bool       `gorm:"default:true"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}