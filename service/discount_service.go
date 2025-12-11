package service

import (
	"context"
	"fmt"
	"Managemenumkm/domain"
	"Managemenumkm/repository"
	"time"
)

type DiscountService interface {
	// CRUD operations
	Create(tokoID uint, request domain.DiscountCreateRequest) (domain.Discount, error)
	Update(request domain.DiscountUpdateRequest) (domain.Discount, error)
	Delete(discountID uint) error
	FindByID(discountID uint) (domain.Discount, error)
	FindByTokoID(tokoID uint) ([]domain.Discount, error)
	FindWithFilter(filter domain.DiscountFilter) ([]domain.Discount, error)

	// Application logic
	GetApplicableDiscounts(ctx context.Context, memberID uint, items []domain.TransactionItem) ([]domain.Discount, error)
	ValidateDiscount(ctx context.Context, discountID uint, memberID uint, cartTotal float64) (*DiscountValidationResult, error)
	ApplyDiscount(ctx context.Context, discount *domain.Discount, subtotal float64, member *domain.Member) (*DiscountApplication, error)
	RecordDiscountUsage(ctx context.Context, discountID uint, transactionID uint, memberID uint, amountUsed float64) error
	GetActiveSpecialOffers(ctx context.Context, tokoID uint) ([]domain.SpecialOffer, error)
	GetActiveFlashSales(ctx context.Context, tokoID uint) ([]domain.FlashSale, error)
}

type DiscountValidationResult struct {
	IsValid      bool
	Discount     *domain.Discount
	ErrorMessage string
	DiscountType string
	MaxDiscount  float64
}

type DiscountApplication struct {
	DiscountAmount float64
	FinalAmount    float64
	Description    string
}

type discountServiceImpl struct {
	discountRepo repository.DiscountRepository
	memberRepo   repository.MemberRepository
}

func NewDiscountService(discountRepo repository.DiscountRepository, memberRepo repository.MemberRepository) DiscountService {
	return &discountServiceImpl{
		discountRepo: discountRepo,
		memberRepo:   memberRepo,
	}
}

func (s *discountServiceImpl) GetApplicableDiscounts(ctx context.Context, memberID uint, items []domain.TransactionItem) ([]domain.Discount, error) {
	var applicableDiscounts []domain.Discount

	// Get all active discounts - this would need implementation in the repository
	// For now, return empty slice
	return applicableDiscounts, nil
}

func (s *discountServiceImpl) ValidateDiscount(ctx context.Context, discountID uint, memberID uint, cartTotal float64) (*DiscountValidationResult, error) {
	// Get discount
	discount, err := s.discountRepo.FindByID(discountID)
	if err != nil {
		return nil, err
	}

	// Check if discount is active and within valid date range
	now := time.Now()
	if !discount.IsActive || now.Before(discount.StartDate) || now.After(discount.EndDate) {
		return &DiscountValidationResult{
			IsValid:      false,
			ErrorMessage: "Discount is not active or expired",
		}, nil
	}

	// Get member if applicable
	var member domain.Member
	if memberID != 0 {
		member, err = s.memberRepo.FindByID(memberID)
		if err != nil {
			return nil, err
		}
	}

	// Check member tier restriction
	if discount.MemberTierID != nil && memberID != 0 {
		if member.MemberTierID != *discount.MemberTierID {
			return &DiscountValidationResult{
				IsValid:      false,
				ErrorMessage: "This discount is not available for your membership tier",
			}, nil
		}
	}

	// Check minimum purchase requirement
	if discount.MinPurchase > 0 && cartTotal < discount.MinPurchase {
		return &DiscountValidationResult{
			IsValid:      false,
			ErrorMessage: "Minimum purchase amount not met",
		}, nil
	}

	// Check usage limits
	if discount.UsageLimit > 0 {
		usage, err := s.discountRepo.GetDiscountUsage(discountID)
		if err != nil {
			return nil, err
		}

		if usage >= discount.UsageLimit {
			return &DiscountValidationResult{
				IsValid:      false,
				ErrorMessage: "Discount usage limit exceeded",
			}, nil
		}
	}

	return &DiscountValidationResult{
		IsValid:      true,
		Discount:     &discount,
		DiscountType: discount.Type,
		MaxDiscount:  discount.MaxDiscount,
	}, nil
}

func (s *discountServiceImpl) ApplyDiscount(ctx context.Context, discount *domain.Discount, subtotal float64, member *domain.Member) (*DiscountApplication, error) {
	var discountAmount float64
	var description string

	switch discount.Type {
	case "PERCENTAGE":
		discountAmount = subtotal * (discount.DiscountValue / 100)
		description = "Percentage discount"

	case "FIXED_AMOUNT":
		discountAmount = discount.DiscountValue
		description = "Fixed amount discount"

	case "BUY_X_GET_Y":
		// This would need more complex logic based on cart items
		discountAmount = discount.DiscountValue
		description = "Buy X Get Y discount"

	case "BULK":
		// This would need quantity-based calculation
		discountAmount = discount.DiscountValue
		description = "Bulk purchase discount"

	default:
		discountAmount = 0
		description = "Unknown discount type"
	}

	// Apply maximum discount limit
	if discount.MaxDiscount > 0 && discountAmount > discount.MaxDiscount {
		discountAmount = discount.MaxDiscount
	}

	// Ensure discount doesn't make amount negative
	if discountAmount > subtotal {
		discountAmount = subtotal
	}

	finalAmount := subtotal - discountAmount

	return &DiscountApplication{
		DiscountAmount: discountAmount,
		FinalAmount:    finalAmount,
		Description:    description,
	}, nil
}

func (s *discountServiceImpl) RecordDiscountUsage(ctx context.Context, discountID uint, transactionID uint, memberID uint, amountUsed float64) error {
	usage := domain.DiscountUsage{
		DiscountID:     discountID,
		TransactionID:  transactionID,
		OriginalAmount: amountUsed, // This would be calculated properly
		DiscountAmount: amountUsed, // This would be calculated properly
		FinalAmount:    amountUsed, // This would be calculated properly
		CreatedAt:      time.Now(),
	}

	// This method would need to be implemented in the repository
	_, err := s.discountRepo.CreateDiscountUsage(usage)
	return err
}

func (s *discountServiceImpl) GetActiveSpecialOffers(ctx context.Context, tokoID uint) ([]domain.SpecialOffer, error) {
	return s.discountRepo.FindActiveSpecialOffers(tokoID)
}

func (s *discountServiceImpl) GetActiveFlashSales(ctx context.Context, tokoID uint) ([]domain.FlashSale, error) {
	return s.discountRepo.FindActiveFlashSales(tokoID)
}

// CRUD operations implementation

func (s *discountServiceImpl) Create(tokoID uint, request domain.DiscountCreateRequest) (domain.Discount, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return domain.Discount{}, err
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return domain.Discount{}, err
	}

	// Validate date range
	if endDate.Before(startDate) {
		return domain.Discount{}, fmt.Errorf("end date must be after start date")
	}

	discount := domain.Discount{
		TokoID:              tokoID,
		Name:                request.Name,
		Description:         request.Description,
		Type:                request.Type,
		DiscountValue:       request.DiscountValue,
		MinPurchase:         request.MinPurchase,
		MaxDiscount:         request.MaxDiscount,
		MinQuantity:         request.MinQuantity,
		BuyQuantity:         request.BuyQuantity,
		GetQuantity:         request.GetQuantity,
		GetDiscountPercent:  request.GetDiscountPercent,
		ApplicableTo:        request.ApplicableTo,
		ProductID:           request.ProductID,
		CategoryID:          request.CategoryID,
		MemberTierID:        request.MemberTierID,
		StartDate:           startDate,
		EndDate:             endDate,
		UsageLimit:          request.UsageLimit,
		IsStackable:         request.IsStackable,
		IsActive:            true,
	}

	return s.discountRepo.Create(discount)
}

func (s *discountServiceImpl) Update(request domain.DiscountUpdateRequest) (domain.Discount, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return domain.Discount{}, err
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return domain.Discount{}, err
	}

	// Validate date range
	if endDate.Before(startDate) {
		return domain.Discount{}, fmt.Errorf("end date must be after start date")
	}

	discount := domain.Discount{
		ID:                  request.ID,
		Name:                request.Name,
		Description:         request.Description,
		Type:                request.Type,
		DiscountValue:       request.DiscountValue,
		MinPurchase:         request.MinPurchase,
		MaxDiscount:         request.MaxDiscount,
		MinQuantity:         request.MinQuantity,
		BuyQuantity:         request.BuyQuantity,
		GetQuantity:         request.GetQuantity,
		GetDiscountPercent:  request.GetDiscountPercent,
		ApplicableTo:        request.ApplicableTo,
		ProductID:           request.ProductID,
		CategoryID:          request.CategoryID,
		MemberTierID:        request.MemberTierID,
		StartDate:           startDate,
		EndDate:             endDate,
		UsageLimit:          request.UsageLimit,
		IsStackable:         request.IsStackable,
		IsActive:            request.IsActive,
	}

	return s.discountRepo.Update(discount)
}

func (s *discountServiceImpl) Delete(discountID uint) error {
	return s.discountRepo.Delete(discountID)
}

func (s *discountServiceImpl) FindByID(discountID uint) (domain.Discount, error) {
	return s.discountRepo.FindByID(discountID)
}

func (s *discountServiceImpl) FindByTokoID(tokoID uint) ([]domain.Discount, error) {
	return s.discountRepo.FindByTokoID(tokoID)
}

func (s *discountServiceImpl) FindWithFilter(filter domain.DiscountFilter) ([]domain.Discount, error) {
	return s.discountRepo.FindWithFilter(filter)
}

// Helper methods

func (s *discountServiceImpl) isDiscountApplicable(discount *domain.Discount, member *domain.Member, items []domain.TransactionItem, cartTotal float64) bool {
	// Check basic validity
	now := time.Now()
	if !discount.IsActive || now.Before(discount.StartDate) || now.After(discount.EndDate) {
		return false
	}

	// Check member tier restriction
	if discount.MemberTierID != nil && member != nil {
		if member.MemberTierID != *discount.MemberTierID {
			return false
		}
	}

	// Check minimum purchase
	if discount.MinPurchase > 0 && cartTotal < discount.MinPurchase {
		return false
	}

	return true
}

func (s *discountServiceImpl) calculateCartTotal(items []domain.TransactionItem) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}