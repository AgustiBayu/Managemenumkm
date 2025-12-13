package service

import (
	"context"
	"Managemenumkm/domain"
	"Managemenumkm/repository"
	"errors"
	"fmt"
	"time"
)

type CheckoutService interface {
	ProcessCheckout(ctx context.Context, request *CheckoutRequest) (*CheckoutResponse, error)
	CalculateOrderSummary(ctx context.Context, items []CheckoutItem, memberID *uint) (*OrderSummary, error)
	ApplyDiscount(ctx context.Context, discountID uint, memberID uint, cartTotal float64) (*DiscountApplication, error)
	RedeemPoints(ctx context.Context, memberID uint, pointsToRedeem int) (*PointsRedemption, error)
}

type CheckoutRequest struct {
	TokoID         uint              `json:"toko_id"`
	UserID         uint              `json:"user_id"`
	MemberID       *uint             `json:"member_id,omitempty"`
	Items          []CheckoutItem    `json:"items"`
	PaymentMethod  string            `json:"payment_method"`
	DiscountIDs    []uint            `json:"discount_ids,omitempty"`
	PointsToRedeem int               `json:"points_to_redeem,omitempty"`
}

type CheckoutItem struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CheckoutResponse struct {
	TransactionID    uint            `json:"transaction_id"`
	OrderSummary     *OrderSummary   `json:"order_summary"`
	MemberInfo       *MemberInfo     `json:"member_info,omitempty"`
	AppliedDiscounts []string        `json:"applied_discounts,omitempty"`
	PointsEarned     int             `json:"points_earned,omitempty"`
	PointsRedeemed   int             `json:"points_redeemed,omitempty"`
	Success          bool            `json:"success"`
	Message          string          `json:"message"`
}

type OrderSummary struct {
	Subtotal        float64            `json:"subtotal"`
	DiscountAmount  float64            `json:"discount_amount"`
	PointsDiscount  float64            `json:"points_discount"`
	TaxAmount       float64            `json:"tax_amount"`
	TotalAmount     float64            `json:"total_amount"`
	Items           []ItemSummary      `json:"items"`
	AppliedDiscounts []DiscountInfo    `json:"applied_discounts,omitempty"`
}

type DiscountInfo struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type ItemSummary struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type MemberInfo struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	MemberCode     string  `json:"member_code"`
	Tier           string  `json:"tier"`
	PointsBalance  int     `json:"points_balance"`
	PointsEarned   int     `json:"points_earned"`
	PointsRedeemed int     `json:"points_redeemed"`
}

type PointsRedemption struct {
	PointsUsed     int     `json:"points_used"`
	DiscountAmount float64 `json:"discount_amount"`
	NewBalance     int     `json:"new_balance"`
}

type checkoutServiceImpl struct {
	transactionRepo  repository.TransactionRepository
	memberRepo       repository.MemberRepository
	productRepo      repository.ProductRepository
	pointsService    PointsCalculationService
	discountService  DiscountService
	taxRate          float64
}

func NewCheckoutService(
	transactionRepo repository.TransactionRepository,
	memberRepo repository.MemberRepository,
	productRepo repository.ProductRepository,
	pointsService PointsCalculationService,
	discountService DiscountService,
) CheckoutService {
	return &checkoutServiceImpl{
		transactionRepo: transactionRepo,
		memberRepo:      memberRepo,
		productRepo:     productRepo,
		pointsService:   pointsService,
		discountService: discountService,
		taxRate:         0.10, // 10% tax rate
	}
}

func (s *checkoutServiceImpl) ProcessCheckout(ctx context.Context, request *CheckoutRequest) (*CheckoutResponse, error) {
	// Calculate order summary
	orderSummary, err := s.calculateOrderSummaryInternal(ctx, request)
	if err != nil {
		return &CheckoutResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to calculate order: %v", err),
		}, err
	}

	// Create transaction
	transaction := domain.Transaction{
		TokoID:          request.TokoID,
		UserID:          request.UserID,
		TotalAmount:     orderSummary.TotalAmount,
		PaymentMethod:   request.PaymentMethod,
		TransactionDate: time.Now(),
	}

	// Convert items to domain items
	var transactionItems []domain.TransactionItem
	for i, item := range request.Items {
		transactionItems = append(transactionItems, domain.TransactionItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})

		// Update item summary with product name if available
		if i < len(orderSummary.Items) {
			orderSummary.Items[i].ProductID = item.ProductID
		}
	}

	// Save transaction (this would need to be implemented properly in the repository)
	// For now, we'll simulate the transaction creation
	transactionID := uint(time.Now().Unix()) // Temporary ID generation

	// Handle member-related operations
	var memberInfo *MemberInfo
	var pointsEarned int

	if request.MemberID != nil {
		member, err := s.memberRepo.FindByID(*request.MemberID)
		if err == nil {
			memberInfo = &MemberInfo{
				ID:             member.ID,
				Name:           member.Name,
				MemberCode:     member.MemberCode,
				Tier: "Standard",
				PointsBalance:  member.TotalPoints,
				PointsRedeemed: request.PointsToRedeem,
			}

			// Calculate and award points
			pointsEarned, err = s.pointsService.CalculatePoints(ctx, &transaction, &member)
			if err == nil && pointsEarned > 0 {
				newBalance := member.TotalPoints + pointsEarned - request.PointsToRedeem
				err = s.memberRepo.UpdatePoints(member.ID, newBalance)
				if err == nil {
					memberInfo.PointsEarned = pointsEarned
					memberInfo.PointsBalance = newBalance
				}
			}
		}
	}

	// Record discount usage
	if len(request.DiscountIDs) > 0 {
		for _, discountID := range request.DiscountIDs {
			err = s.discountService.RecordDiscountUsage(ctx, discountID, transactionID, *request.MemberID, orderSummary.DiscountAmount)
			if err != nil {
				// Log error but don't fail the transaction
				fmt.Printf("Failed to record discount usage: %v\n", err)
			}
		}
	}

	return &CheckoutResponse{
		TransactionID:    transactionID,
		OrderSummary:     orderSummary,
		MemberInfo:       memberInfo,
		AppliedDiscounts: s.getAppliedDiscountNames(orderSummary.AppliedDiscounts),
		PointsEarned:     pointsEarned,
		PointsRedeemed:   request.PointsToRedeem,
		Success:          true,
		Message:          "Transaction completed successfully",
	}, nil
}

func (s *checkoutServiceImpl) CalculateOrderSummary(ctx context.Context, items []CheckoutItem, memberID *uint) (*OrderSummary, error) {
	request := &CheckoutRequest{
		Items:    items,
		MemberID: memberID,
	}
	return s.calculateOrderSummaryInternal(ctx, request)
}

func (s *checkoutServiceImpl) calculateOrderSummaryInternal(ctx context.Context, request *CheckoutRequest) (*OrderSummary, error) {
	// Calculate subtotal
	var itemSummaries []ItemSummary
	subtotal := 0.0

	for _, item := range request.Items {
		itemSubtotal := float64(item.Quantity) * item.Price
		subtotal += itemSubtotal

		itemSummaries = append(itemSummaries, ItemSummary{
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: itemSubtotal,
		})
	}

	// Calculate total before discounts
	totalBeforeDiscounts := subtotal

	// Apply discounts
	var appliedDiscounts []DiscountInfo
	discountAmount := 0.0

	if len(request.DiscountIDs) > 0 && request.MemberID != nil {
		for _, discountID := range request.DiscountIDs {
			discountApp, err := s.ApplyDiscount(ctx, discountID, *request.MemberID, subtotal)
			if err == nil {
				discountAmount += discountApp.DiscountAmount
				appliedDiscounts = append(appliedDiscounts, DiscountInfo{
					ID:          discountID,
					Name:        "Applied Discount", // Would fetch from discount
					Type:        discountApp.Description,
					Amount:      discountApp.DiscountAmount,
					Description: discountApp.Description,
				})
			}
		}
	}

	// Apply points redemption
	pointsDiscount := 0.0
	if request.PointsToRedeem > 0 && request.MemberID != nil {
		// Check if member exists before redeeming points
		_, err := s.memberRepo.FindByID(*request.MemberID)
		if err == nil {
			pointsValue, err := s.pointsService.CalculatePointsValue(ctx, request.PointsToRedeem, "Standard")
			if err == nil {
				pointsDiscount = pointsValue
			}
		}
	}

	// Calculate tax (applied after discounts)
	subtotalAfterDiscounts := totalBeforeDiscounts - discountAmount - pointsDiscount
	taxAmount := subtotalAfterDiscounts * s.taxRate

	// Final total
	totalAmount := subtotalAfterDiscounts + taxAmount

	return &OrderSummary{
		Subtotal:         subtotal,
		DiscountAmount:   discountAmount,
		PointsDiscount:   pointsDiscount,
		TaxAmount:        taxAmount,
		TotalAmount:      totalAmount,
		Items:            itemSummaries,
		AppliedDiscounts: appliedDiscounts,
	}, nil
}

func (s *checkoutServiceImpl) ApplyDiscount(ctx context.Context, discountID uint, memberID uint, cartTotal float64) (*DiscountApplication, error) {
	// Get member info
	member, err := s.memberRepo.FindByID(memberID)
	if err != nil {
		return nil, err
	}

	// Get discount
	// This would need to be implemented in the discount repository
	// For now, we'll use the validation service
	validation, err := s.discountService.ValidateDiscount(ctx, discountID, memberID, cartTotal)
	if err != nil {
		return nil, err
	}

	if !validation.IsValid {
		return nil, errors.New(validation.ErrorMessage)
	}

	// Apply discount
	return s.discountService.ApplyDiscount(ctx, validation.Discount, cartTotal, &member)
}

func (s *checkoutServiceImpl) RedeemPoints(ctx context.Context, memberID uint, pointsToRedeem int) (*PointsRedemption, error) {
	// Process points redemption
	_, err := s.pointsService.ProcessPointsRedemption(ctx, memberID, pointsToRedeem)
	if err != nil {
		return nil, err
	}

	// Get updated member balance
	member, err := s.memberRepo.FindByID(memberID)
	if err != nil {
		return nil, err
	}

	// Calculate discount amount
	discountAmount, err := s.pointsService.CalculatePointsValue(ctx, pointsToRedeem, "Standard")
	if err != nil {
		return nil, err
	}

	return &PointsRedemption{
		PointsUsed:     pointsToRedeem,
		DiscountAmount: discountAmount,
		NewBalance:     member.TotalPoints,
	}, nil
}

// Helper methods

func (s *checkoutServiceImpl) getAppliedDiscountNames(discounts []DiscountInfo) []string {
	var names []string
	for _, discount := range discounts {
		names = append(names, discount.Name)
	}
	return names
}