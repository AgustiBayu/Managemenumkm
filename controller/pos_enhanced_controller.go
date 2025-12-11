package controller

import (
	"Managemenumkm/helper"
	"Managemenumkm/middleware"
	"Managemenumkm/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PosEnhancedController interface {
	ShowPosPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	SearchMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetMemberPoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ApplyDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAvailableDiscounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	RedeemPoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ProcessCheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type PosEnhancedControllerImpl struct {
	CheckoutService   service.CheckoutService
	MemberService     service.MemberService
	ProductService    service.ProductService
	DiscountService   service.DiscountService
	PointsService     service.PointsCalculationService
}

func NewPosEnhancedController(
	checkoutService service.CheckoutService,
	memberService service.MemberService,
	productService service.ProductService,
	discountService service.DiscountService,
	pointsService service.PointsCalculationService,
) PosEnhancedController {
	return &PosEnhancedControllerImpl{
		CheckoutService: checkoutService,
		MemberService:   memberService,
		ProductService:  productService,
		DiscountService: discountService,
		PointsService:   pointsService,
	}
}

func (controller *PosEnhancedControllerImpl) ShowPosPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	helper.RenderTemplate(w, "template/pos/pos.html", map[string]interface{}{
		"title":  "Enhanced Point of Sale",
		"UserID": claims.UserID,
		"TokoID": claims.TokoID,
	})
}

func (controller *PosEnhancedControllerImpl) SearchMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)

	// Get search query
	search := r.URL.Query().Get("search")
	if search == "" {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Search query is required")
		return
	}

	// Search member by code, phone, or name
	member, err := controller.MemberService.FindByMemberCode(search, claims.TokoID)
	if err != nil {
		// Try searching by phone
		member, err = controller.MemberService.FindByPhone(search, claims.TokoID)
		if err != nil {
			helper.WriteErrorResponse(w, http.StatusNotFound, "Member not found")
			return
		}
	}

	// Prepare response
	response := map[string]interface{}{
		"id":             member.ID,
		"member_code":    member.MemberCode,
		"name":           member.Name,
		"email":          member.Email,
		"phone":          member.Phone,
		"tier":           member.MemberTier.Name,
		"points_balance": member.TotalPoints,
		"total_spent":    member.TotalSpent,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *PosEnhancedControllerImpl) GetMemberPoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = r.Context().Value(middleware.ClaimsKey).(*helper.Claims) // Get claims for authorization check if needed

	memberIDStr := ps.ByName("id")
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	member, err := controller.MemberService.GetMemberByID(uint(memberID))
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusNotFound, "Member not found")
		return
	}

	// Calculate points value
	pointsValue, err := controller.PointsService.CalculatePointsValue(r.Context(), member.TotalPoints, member.MemberTier.Name)
	if err != nil {
		pointsValue = 0
	}

	response := map[string]interface{}{
		"member_id":      member.ID,
		"member_code":    member.MemberCode,
		"name":           member.Name,
		"tier":           member.MemberTier.Name,
		"points_balance": member.TotalPoints,
		"points_value":   pointsValue,
		"total_spent":    member.TotalSpent,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *PosEnhancedControllerImpl) ApplyDiscount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request struct {
		MemberID  uint    `json:"member_id"`
		DiscountID uint   `json:"discount_id"`
		CartTotal float64 `json:"cart_total"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate and apply discount
	discountApp, err := controller.CheckoutService.ApplyDiscount(r.Context(), request.DiscountID, request.MemberID, request.CartTotal)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"discount_amount": discountApp.DiscountAmount,
		"final_amount":    discountApp.FinalAmount,
		"description":     discountApp.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *PosEnhancedControllerImpl) GetAvailableDiscounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = r.URL.Query().Get("member_id") // Get member_id but don't use it for now

	// This would need to be implemented in the discount service
	// For now, return empty array
	discounts := []interface{}{}

	response := map[string]interface{}{
		"discounts": discounts,
		"count":     len(discounts),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *PosEnhancedControllerImpl) RedeemPoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request struct {
		MemberID       uint `json:"member_id"`
		PointsToRedeem int  `json:"points_to_redeem"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Process points redemption
	redemption, err := controller.CheckoutService.RedeemPoints(r.Context(), request.MemberID, request.PointsToRedeem)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"points_used":     redemption.PointsUsed,
		"discount_amount": redemption.DiscountAmount,
		"new_balance":     redemption.NewBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *PosEnhancedControllerImpl) ProcessCheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)

	var request service.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set toko ID from claims
	request.TokoID = claims.TokoID
	request.UserID = claims.UserID

	// Process checkout
	response, err := controller.CheckoutService.ProcessCheckout(r.Context(), &request)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}