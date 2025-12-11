package repository

import (
	"Managemenumkm/domain"
	"time"
)

type DiscountRepository interface {
	// Basic CRUD operations
	Create(discount domain.Discount) (domain.Discount, error)
	Update(discount domain.Discount) (domain.Discount, error)
	Delete(discountID uint) error
	FindByID(discountID uint) (domain.Discount, error)
	FindByTokoID(tokoID uint) ([]domain.Discount, error)

	// Advanced queries
	FindActiveDiscounts(tokoID uint) ([]domain.Discount, error)
	FindValidDiscounts(tokoID uint, currentDate time.Time) ([]domain.Discount, error)
	FindWithFilter(filter domain.DiscountFilter) ([]domain.Discount, error)

	// Discount usage tracking
	GetDiscountUsage(discountID uint) (int, error)
	IncrementUsage(discountID uint) error

	// Special offers
	CreateSpecialOffer(specialOffer domain.SpecialOffer) (domain.SpecialOffer, error)
	UpdateSpecialOffer(specialOffer domain.SpecialOffer) (domain.SpecialOffer, error)
	DeleteSpecialOffer(specialOfferID uint) error
	FindActiveSpecialOffers(tokoID uint) ([]domain.SpecialOffer, error)
	FindSpecialOfferProducts(specialOfferID uint) ([]domain.SpecialOfferProduct, error)

	// Flash sales
	CreateFlashSale(flashSale domain.FlashSale) (domain.FlashSale, error)
	UpdateFlashSale(flashSale domain.FlashSale) (domain.FlashSale, error)
	DeleteFlashSale(flashSaleID uint) error
	FindActiveFlashSales(tokoID uint) ([]domain.FlashSale, error)
	FindFlashSaleProducts(flashSaleID uint) ([]domain.FlashSaleProduct, error)

	// Discount usage tracking
	CreateDiscountUsage(usage domain.DiscountUsage) (domain.DiscountUsage, error)
	FindDiscountUsagesByTransaction(transactionID uint) ([]domain.DiscountUsage, error)
	FindProductDiscounts(productID uint, tokoID uint) ([]domain.Discount, error)
	FindCategoryDiscounts(categoryID uint, tokoID uint) ([]domain.Discount, error)
	FindMemberTierDiscounts(memberTierID uint, tokoID uint) ([]domain.Discount, error)
}