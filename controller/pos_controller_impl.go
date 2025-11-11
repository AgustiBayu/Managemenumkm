package controller

import (
	"Managemenumkm/helper"
	"Managemenumkm/middleware"
	"Managemenumkm/service"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PosControllerImpl struct {
	TransactionService service.TransactionService
	ProductService     service.ProductService
}

func NewPosController(transactionService service.TransactionService, productService service.ProductService) PosController {
	return &PosControllerImpl{
		TransactionService: transactionService,
		ProductService:     productService,
	}
}

func (controller *PosControllerImpl) ShowPosPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve claims from context, set by the middleware
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*helper.Claims)
	if !ok || claims == nil {
		// If claims are not found in context, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Pass user and toko ID from claims to the template
	helper.RenderTemplate(w, "template/pos/pos.html", map[string]interface{}{
		"title":  "Point of Sale",
		"UserID": claims.UserID,
		"TokoID": claims.TokoID,
	})
}

func (controller *PosControllerImpl) GetProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	products, err := controller.ProductService.FindAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (controller *PosControllerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var checkoutRequest helper.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&checkoutRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Note: In a real app, UserID and TokoID should be retrieved from the authenticated session (e.g., JWT parsed in middleware).
	// For now, we trust the client, but this should be improved for security.

	transaction, err := controller.TransactionService.Checkout(checkoutRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
