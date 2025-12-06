package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/service"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type MemberController interface {
	CreateMember(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	UpdateMember(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	DeleteMember(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetMemberByID(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetMembersByTokoID(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindMemberByCode(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindMemberByPhone(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetMemberTransactions(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetPointsBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	RedeemPoints(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetMemberTiers(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	ImportMembers(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	// Web Interface Methods
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetMembersAPI(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type MemberControllerImpl struct {
	MemberService service.MemberService
}

func NewMemberController(memberService service.MemberService) MemberController {
	return &MemberControllerImpl{
		MemberService: memberService,
	}
}

func (controller *MemberControllerImpl) CreateMember(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var request domain.Member
		err := decoder.Decode(&request)
		helper.PanicIfError(err)

		// For now, use a default tokoID or get from session
		// TODO: Get tokoID from context or session
		request.TokoID = 1 // Default tokoID for now

		member, err := controller.MemberService.CreateMember(request)
		if err != nil {
			http.Error(w, "Failed to create member: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) UpdateMember(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodPut || r.Method == http.MethodPost {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request domain.Member
		err = decoder.Decode(&request)
		helper.PanicIfError(err)

		// Debug logging to help troubleshoot tier update issues
		fmt.Printf("UpdateMember request for ID %d: %+v\n", memberID, request)
		fmt.Printf("MemberTierID in request: %d\n", request.MemberTierID)

		member, err := controller.MemberService.UpdateMember(uint(memberID), request)
		if err != nil {
			fmt.Printf("Error updating member: %v\n", err)
			http.Error(w, "Failed to update member: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Debug logging for successful update
		fmt.Printf("Successfully updated member: %+v\n", member)
		fmt.Printf("Updated MemberTierID: %d\n", member.MemberTierID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) DeleteMember(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodDelete {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = controller.MemberService.DeleteMember(uint(memberID))
		if err != nil {
			http.Error(w, "Failed to delete member: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Member deleted successfully"})
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) GetMemberByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		member, err := controller.MemberService.GetMemberByID(uint(memberID))
		if err != nil {
			http.Error(w, "Member not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) GetMembersByTokoID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		// For now, use a default tokoID
		// TODO: Get tokoID from context or session
		tokoID := uint(1) // Default tokoID for now

		members, err := controller.MemberService.GetMembersByTokoID(tokoID)
		if err != nil {
			http.Error(w, "Failed to get members: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(members)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) FindMemberByCode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		// TODO: Get tokoID from context or session
		tokoID := uint(1) // Default tokoID for now

		// Try to find by member code first
		member, err := controller.MemberService.FindByMemberCode(query, tokoID)
		if err == nil {
			// Found by code, return as single result array
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]domain.Member{member})
			return
		}

		// Try to find by phone if code search failed
		member, err = controller.MemberService.FindByPhone(query, tokoID)
		if err == nil {
			// Found by phone, return as single result array
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]domain.Member{member})
			return
		}

		// If both searches failed, return empty array
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]domain.Member{})
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// API endpoint to get all members for frontend
func (controller *MemberControllerImpl) GetMembersAPI(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		tokoID := uint(1) // Default tokoID for now

		// Get filter parameters from query string
		search := r.URL.Query().Get("search")
		tier := r.URL.Query().Get("tier")
		status := r.URL.Query().Get("status")

		// Create filter criteria
		filter := domain.MemberFilter{
			TokoID: tokoID,
			Search: search,
			Tier:   tier,
			Status: status,
		}

		members, err := controller.MemberService.GetMembersWithFilter(filter)
		if err != nil {
			http.Error(w, "Failed to get members: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(members)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) FindMemberByPhone(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		phone := params.ByName("phone")
		if phone == "" {
			http.Error(w, "Phone number is required", http.StatusBadRequest)
			return
		}

		// TODO: Get tokoID from context or session
		tokoID := uint(1) // Default tokoID for now

		member, err := controller.MemberService.FindByPhone(phone, tokoID)
		if err != nil {
			http.Error(w, "Member not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) GetMemberTransactions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 20
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		transactions, err := controller.MemberService.GetMemberTransactions(uint(memberID), limit, offset)
		if err != nil {
			http.Error(w, "Failed to get transactions: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) GetPointsBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		balance, err := controller.MemberService.GetPointsBalance(uint(memberID))
		if err != nil {
			http.Error(w, "Failed to get points balance: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"points_balance": balance})
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

type RedeemPointsRequest struct {
	Points      int    `json:"points" validate:"required,min=1"`
	Description string `json:"description" validate:"required"`
}

func (controller *MemberControllerImpl) RedeemPoints(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodPost {
		idParam := params.ByName("id")
		if idParam == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}
		memberID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request RedeemPointsRequest
		err = decoder.Decode(&request)
		helper.PanicIfError(err)

		transaction, err := controller.MemberService.RedeemPoints(uint(memberID), request.Points, request.Description)
		if err != nil {
			http.Error(w, "Failed to redeem points: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transaction)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) GetMemberTiers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		tokoID := uint(1) // Default tokoID for now

		memberTiers, err := controller.MemberService.GetMemberTiers(tokoID)
		if err != nil {
			http.Error(w, "Failed to get member tiers: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(memberTiers)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// ImportMembersRequest represents the request body for importing members
type ImportMembersRequest struct {
	Members        []service.ImportMember `json:"members"`
	SkipDuplicates bool                  `json:"skipDuplicates"`
	UpdateExisting bool                  `json:"updateExisting"`
}

func (controller *MemberControllerImpl) ImportMembers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request ImportMembersRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if len(request.Members) == 0 {
		http.Error(w, "No members to import", http.StatusBadRequest)
		return
	}

	// For now, use a default tokoID
	// TODO: Get tokoID from context or session
	tokoID := uint(1)

	response, err := controller.MemberService.ImportMembers(request.Members, tokoID, request.SkipDuplicates, request.UpdateExisting)
	if err != nil {
		http.Error(w, "Failed to import members: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Template functions
func (controller *MemberControllerImpl) addTemplateFunctions(tmpl *template.Template) *template.Template {
	funcMap := template.FuncMap{
		"currencyFormat": func(amount float64) string {
			if amount == 0 {
				return "Rp 0"
			}
			return fmt.Sprintf("Rp %,.0f", amount)
		},
		"formatDate": func(date time.Time) string {
			if date.IsZero() {
				return "-"
			}
			return date.Format("02 Jan 2006")
		},
		"substr": func(s string, start, length int) string {
			if len(s) <= start {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	return tmpl.Funcs(funcMap)
}

// Web Interface Implementation
func (controller *MemberControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Add timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	r = r.WithContext(ctx)

	// Get member data
	tokoID := uint(1) // Default tokoID for now

	// Remove AutoInactiveMembers from request-time to avoid blocking
	// This should be run in background/cron, not during page load

	members, err := controller.MemberService.GetMembersByTokoID(tokoID)
	if err != nil {
		http.Error(w, "Failed to get members: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if context was cancelled
	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
		// Continue
	}

	// Get member tiers
	memberTiers, err := controller.MemberService.GetMemberTiers(tokoID)
	if err != nil {
		// Continue without tiers if error
		memberTiers = []domain.MemberTier{}
	}

	// Prepare member data for frontend (add calculated fields)
	// Limit to first 50 members to avoid page loading issues
	displayLimit := 50
	if len(members) > displayLimit {
		members = members[:displayLimit]
	}

	memberData := make([]map[string]interface{}, len(members))
	for i, member := range members {
		memberData[i] = map[string]interface{}{
			"id":                member.ID,
			"memberCode":        member.MemberCode,
			"name":              member.Name,
			"phone":             member.Phone,
			"email":             member.Email,
			"address":           member.Address,
			"birthday":          member.Birthday,
			"gender":            member.Gender,
			"notes":             member.Notes,
			"totalPoints":       member.TotalPoints,
			"totalSpent":        member.TotalSpent,
			"totalTransactions": member.TotalTransactions,
			"status":            member.Status,
			"isActive":          member.IsActive,
			"joinedDate":        member.JoinedDate,
			"lastVisitDate":     member.LastVisitDate,
			"memberTier":        member.MemberTier,
		}
	}

	// Prepare data for template
	data := map[string]interface{}{
		"title":   "Member Management",
		"members": memberData,
		"tiers":   memberTiers,
		"page":    "members",
	}

	// Use helper.RenderTemplate which has proper function mapping
	helper.RenderTemplate(w, "template/member/member_list.html", data)
}

func (controller *MemberControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodGet {
		// Return JSON with member tiers for the form
		// In a real app, you'd render the member add template
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Member create form"})
		return
	}

	if r.Method == http.MethodPost {
		// Handle form data instead of JSON
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		address := r.FormValue("address")
		birthday := r.FormValue("birthday")
		gender := r.FormValue("gender")
		notes := r.FormValue("notes")
		tierId := r.FormValue("tierId")

		if name == "" || phone == "" {
			http.Error(w, "Name and phone are required", http.StatusBadRequest)
			return
		}

		member := domain.Member{
			Name:    name,
			Phone:   phone,
			Email:   email,
			Address: address,
			TokoID:  1, // Default tokoID
		}

		// Parse optional fields
		if birthday != "" {
			member.Birthday = birthday
		}
		if gender != "" {
			member.Gender = gender
		}
		if notes != "" {
			member.Notes = notes
		}
		if tierId != "" {
			if id, err := strconv.ParseUint(tierId, 10, 32); err == nil {
				member.MemberTierID = uint(id)
			}
		}

		savedMember, err := controller.MemberService.CreateMember(member)
		if err != nil {
			http.Error(w, "Failed to create member: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(savedMember)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (controller *MemberControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idParam := params.ByName("memberId")
	if idParam == "" {
		http.Error(w, "Member ID is required", http.StatusBadRequest)
		return
	}
	memberID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid member ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		// Return member data for edit form
		member, err := controller.MemberService.GetMemberByID(uint(memberID))
		if err != nil {
			http.Error(w, "Member not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
		return
	}

	if r.Method == http.MethodPost {
		// Handle form data for update
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		address := r.FormValue("address")
		birthday := r.FormValue("birthday")
		gender := r.FormValue("gender")
		notes := r.FormValue("notes")
		tierId := r.FormValue("tierId")

		if name == "" || phone == "" {
			http.Error(w, "Name and phone are required", http.StatusBadRequest)
			return
		}

		member := domain.Member{
			Name:    name,
			Phone:   phone,
			Email:   email,
			Address: address,
		}

		// Parse optional fields
		if birthday != "" {
			member.Birthday = birthday
		}
		if gender != "" {
			member.Gender = gender
		}
		if notes != "" {
			member.Notes = notes
		}
		if tierId != "" {
			if id, err := strconv.ParseUint(tierId, 10, 32); err == nil {
				member.MemberTierID = uint(id)
			}
		}

		updatedMember, err := controller.MemberService.UpdateMember(uint(memberID), member)
		if err != nil {
			http.Error(w, "Failed to update member: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedMember)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
