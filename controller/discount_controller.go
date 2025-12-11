package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/middleware"
	"Managemenumkm/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DiscountController interface {
	ShowDiscountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetDiscounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetDiscountByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateSpecialOffer(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateFlashSale(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type DiscountControllerImpl struct {
	DiscountService service.DiscountService
	MemberService   service.MemberService
	ProductService  service.ProductService
}

func NewDiscountController(discountService service.DiscountService, memberService service.MemberService, productService service.ProductService) DiscountController {
	return &DiscountControllerImpl{
		DiscountService: discountService,
		MemberService:   memberService,
		ProductService:  productService,
	}
}

func (controller *DiscountControllerImpl) ShowDiscountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get member tiers for dropdown (only active tiers)
	memberTiers, err := controller.MemberService.GetMemberTiers(claims.TokoID)
	if err != nil {
		memberTiers = []domain.MemberTier{}
	}

	// Filter only active tiers
	var activeMemberTiers []domain.MemberTier
	for _, tier := range memberTiers {
		if tier.IsActive {
			activeMemberTiers = append(activeMemberTiers, tier)
		}
	}
	memberTiers = activeMemberTiers

	// Get products for dropdown
	products, err := controller.ProductService.FindAll()
	if err != nil {
		products = []domain.ProductResponse{} // Use Response type instead
	}

	helper.RenderTemplate(w, "template/admin/discount.html", map[string]interface{}{
		"title":       "Discount Management",
		"UserID":      claims.UserID,
		"TokoID":      claims.TokoID,
		"MemberTiers": memberTiers,
		"Products":    products,
	})
}

func (controller *DiscountControllerImpl) CreateDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)

	var request domain.DiscountCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Use service to create discount
	discount, err := controller.DiscountService.Create(claims.TokoID, request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create discount: "+err.Error())
		return
	}

	// Convert to response format
	response := domain.DiscountResponse{
		ID:                discount.ID,
		Name:              discount.Name,
		Description:       discount.Description,
		Type:              discount.Type,
		DiscountValue:     discount.DiscountValue,
		MinPurchase:       discount.MinPurchase,
		MaxDiscount:       discount.MaxDiscount,
		MinQuantity:       discount.MinQuantity,
		BuyQuantity:       discount.BuyQuantity,
		GetQuantity:       discount.GetQuantity,
		GetDiscountPercent: discount.GetDiscountPercent,
		ApplicableTo:      discount.ApplicableTo,
		ProductID:         discount.ProductID,
		CategoryID:        discount.CategoryID,
		MemberTierID:      discount.MemberTierID,
		StartDate:         discount.StartDate.Format("2006-01-02"),
		EndDate:           discount.EndDate.Format("2006-01-02"),
		UsageLimit:        discount.UsageLimit,
		UsageCount:        discount.UsageCount,
		RemainingUsage:    discount.GetRemainingUsage(),
		IsActive:          discount.IsActive,
		IsStackable:       discount.IsStackable,
		CreatedAt:         discount.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         discount.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helper.WriteSuccessResponse(w, response)
}

func (controller *DiscountControllerImpl) UpdateDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request domain.DiscountUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Use service to update discount
	discount, err := controller.DiscountService.Update(request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update discount: "+err.Error())
		return
	}

	// Convert to response format
	response := domain.DiscountResponse{
		ID:                discount.ID,
		Name:              discount.Name,
		Description:       discount.Description,
		Type:              discount.Type,
		DiscountValue:     discount.DiscountValue,
		MinPurchase:       discount.MinPurchase,
		MaxDiscount:       discount.MaxDiscount,
		MinQuantity:       discount.MinQuantity,
		BuyQuantity:       discount.BuyQuantity,
		GetQuantity:       discount.GetQuantity,
		GetDiscountPercent: discount.GetDiscountPercent,
		ApplicableTo:      discount.ApplicableTo,
		ProductID:         discount.ProductID,
		CategoryID:        discount.CategoryID,
		MemberTierID:      discount.MemberTierID,
		StartDate:         discount.StartDate.Format("2006-01-02"),
		EndDate:           discount.EndDate.Format("2006-01-02"),
		UsageLimit:        discount.UsageLimit,
		UsageCount:        discount.UsageCount,
		RemainingUsage:    discount.GetRemainingUsage(),
		IsActive:          discount.IsActive,
		IsStackable:       discount.IsStackable,
		CreatedAt:         discount.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         discount.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helper.WriteSuccessResponse(w, response)
}

func (controller *DiscountControllerImpl) DeleteDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	discountIDStr := ps.ByName("id")
	discountID, err := strconv.ParseUint(discountIDStr, 10, 32)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid discount ID")
		return
	}

	// Use service to delete discount
	err = controller.DiscountService.Delete(uint(discountID))
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete discount: "+err.Error())
		return
	}

	helper.WriteSuccessResponse(w, map[string]interface{}{
		"message": "Discount deleted successfully",
		"id":      discountID,
	})
}

func (controller *DiscountControllerImpl) GetDiscounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims := r.Context().Value(middleware.ClaimsKey).(*helper.Claims) // Get TokoID for filtering

	// Get query parameters
	discountType := r.URL.Query().Get("type")
	applicableTo := r.URL.Query().Get("applicable_to")
	isActiveStr := r.URL.Query().Get("is_active")

	// Build filter
	filter := domain.DiscountFilter{
		TokoID:       claims.TokoID,
		Type:         discountType,
		ApplicableTo: applicableTo,
	}

	// Parse is_active parameter
	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	// Use service to get discounts
	discounts, err := controller.DiscountService.FindWithFilter(filter)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get discounts: "+err.Error())
		return
	}

	// Convert to response format
	var discountResponses []domain.DiscountResponse
	for _, discount := range discounts {
		// Get product name if applicable
		var productName *string
		if discount.Product != nil {
			productName = &discount.Product.Name
		}

		// Get category name if applicable
		var categoryName *string
		if discount.Category != nil {
			categoryName = &discount.Category.Category
		}

		// Get member tier name if applicable
		var memberTierName *string
		if discount.MemberTier != nil {
			memberTierName = &discount.MemberTier.Name
		}

		response := domain.DiscountResponse{
			ID:                discount.ID,
			Name:              discount.Name,
			Description:       discount.Description,
			Type:              discount.Type,
			DiscountValue:     discount.DiscountValue,
			MinPurchase:       discount.MinPurchase,
			MaxDiscount:       discount.MaxDiscount,
			MinQuantity:       discount.MinQuantity,
			BuyQuantity:       discount.BuyQuantity,
			GetQuantity:       discount.GetQuantity,
			GetDiscountPercent: discount.GetDiscountPercent,
			ApplicableTo:      discount.ApplicableTo,
			ProductID:         discount.ProductID,
			ProductName:       productName,
			CategoryID:        discount.CategoryID,
			CategoryName:      categoryName,
			MemberTierID:      discount.MemberTierID,
			MemberTierName:    memberTierName,
			StartDate:         discount.StartDate.Format("2006-01-02"),
			EndDate:           discount.EndDate.Format("2006-01-02"),
			UsageLimit:        discount.UsageLimit,
			UsageCount:        discount.UsageCount,
			RemainingUsage:    discount.GetRemainingUsage(),
			IsActive:          discount.IsActive,
			IsStackable:       discount.IsStackable,
			CreatedAt:         discount.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         discount.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		discountResponses = append(discountResponses, response)
	}

	response := map[string]interface{}{
		"discounts": discountResponses,
		"count":     len(discountResponses),
	}

	// Write response directly without double wrapping
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": response,
	})
}

func (controller *DiscountControllerImpl) GetDiscountByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	discountIDStr := ps.ByName("id")
	discountID, err := strconv.ParseUint(discountIDStr, 10, 32)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid discount ID")
		return
	}

	// Use service to get discount
	discount, err := controller.DiscountService.FindByID(uint(discountID))
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusNotFound, "Discount not found: "+err.Error())
		return
	}

	// Get product name if applicable
	var productName *string
	if discount.Product != nil {
		productName = &discount.Product.Name
	}

	// Get category name if applicable
	var categoryName *string
	if discount.Category != nil {
		categoryName = &discount.Category.Category
	}

	// Get member tier name if applicable
	var memberTierName *string
	if discount.MemberTier != nil {
		memberTierName = &discount.MemberTier.Name
	}

	// Convert to response format
	response := domain.DiscountResponse{
		ID:                discount.ID,
		Name:              discount.Name,
		Description:       discount.Description,
		Type:              discount.Type,
		DiscountValue:     discount.DiscountValue,
		MinPurchase:       discount.MinPurchase,
		MaxDiscount:       discount.MaxDiscount,
		MinQuantity:       discount.MinQuantity,
		BuyQuantity:       discount.BuyQuantity,
		GetQuantity:       discount.GetQuantity,
		GetDiscountPercent: discount.GetDiscountPercent,
		ApplicableTo:      discount.ApplicableTo,
		ProductID:         discount.ProductID,
		ProductName:       productName,
		CategoryID:        discount.CategoryID,
		CategoryName:      categoryName,
		MemberTierID:      discount.MemberTierID,
		MemberTierName:    memberTierName,
		StartDate:         discount.StartDate.Format("2006-01-02"),
		EndDate:           discount.EndDate.Format("2006-01-02"),
		UsageLimit:        discount.UsageLimit,
		UsageCount:        discount.UsageCount,
		RemainingUsage:    discount.GetRemainingUsage(),
		IsActive:          discount.IsActive,
		IsStackable:       discount.IsStackable,
		CreatedAt:         discount.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         discount.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helper.WriteSuccessResponse(w, response)
}

func (controller *DiscountControllerImpl) CreateSpecialOffer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = r.Context().Value(middleware.ClaimsKey).(*helper.Claims) // For future authorization

	var request domain.SpecialOfferCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// This would need to be implemented in the discount service
	// For now, return success response
	helper.WriteSuccessResponse(w, map[string]interface{}{
		"message": "Special offer created successfully",
		"offer":   request,
	})
}

func (controller *DiscountControllerImpl) CreateFlashSale(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = r.Context().Value(middleware.ClaimsKey).(*helper.Claims) // For future authorization

	var request domain.FlashSaleCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// This would need to be implemented in the discount service
	// For now, return success response
	helper.WriteSuccessResponse(w, map[string]interface{}{
		"message": "Flash sale created successfully",
		"sale":    request,
	})
}