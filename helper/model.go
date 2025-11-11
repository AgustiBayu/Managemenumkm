package helper

import (
	"Managemenumkm/domain"
	"time"
)

// ToProductResponse converts a Product domain model to a ProductResponse DTO.
// It calculates the total stock from all associated batches.
func ToProductResponse(product domain.Product) domain.ProductResponse {
	var totalStock uint
	for _, batch := range product.Batches {
		totalStock += batch.Stock
	}

	var categoryResponse domain.ProductCategoryResponse
	if product.Category.ID != 0 {
		categoryResponse = domain.ProductCategoryResponse{
			ID:       product.Category.ID,
			Category: product.Category.Category,
		}
	}

	return domain.ProductResponse{
		ID:              product.ID,
		Name:            product.Name,
		SKU:             product.SKU,
		Price:           product.Price,
		ImageURL:        product.ImageURL,
		TotalStock:      totalStock,
		CategoryID:      product.CategoryID,
		ProductCategory: categoryResponse,
	}
}

// ToProductResponses converts a slice of Product domain models to a slice of ProductResponse DTOs.
func ToProductResponses(products []domain.Product) []domain.ProductResponse {
	var productResponses []domain.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return productResponses
}

// ToProductBatchResponse converts a ProductBatch domain model to a ProductBatchResponse DTO.
func ToProductBatchResponse(batch domain.ProductBatch) domain.ProductBatchResponse {
	return domain.ProductBatchResponse{
		ID:        batch.ID,
		ProductID: batch.ProductID,
		Barcode:   batch.Barcode,
		Exp:       FormatDate(batch.Exp),
		Stock:     batch.Stock,
		DateAdded: FormatDate(batch.DateAdded),
	}
}

// ToProductBatchResponses converts a slice of ProductBatch domain models to a slice of ProductBatchResponse DTOs.
func ToProductBatchResponses(batches []domain.ProductBatch) []domain.ProductBatchResponse {
	var batchResponses []domain.ProductBatchResponse
	for _, batch := range batches {
		batchResponses = append(batchResponses, ToProductBatchResponse(batch))
	}
	return batchResponses
}

// ToProductBatch converts a ProductBatchUpdateRequest DTO to a ProductBatch domain model.
func ToProductBatch(request domain.ProductBatchUpdateRequest) (domain.ProductBatch, error) {
	expTime, err := time.Parse("2006-01-02", request.Exp)
	if err != nil {
		return domain.ProductBatch{}, err
	}
	return domain.ProductBatch{
		ID:        request.ID,
		ProductID: request.ProductID,
		Stock:     request.Stock,
		Exp:       expTime,
	}, nil
}

func ToProductCategoryResponse(category *domain.ProductCategory) *domain.ProductCategoryResponse {
	return &domain.ProductCategoryResponse{
		ID:       category.ID,
		Category: category.Category,
	}
}
func ToProductCategoryResponses(categories []*domain.ProductCategory) []*domain.ProductCategoryResponse {
	var responses []*domain.ProductCategoryResponse
	for _, category := range categories {
		responses = append(responses, ToProductCategoryResponse(category))
	}
	return responses
}

func ToUserResponse(user *domain.User) *domain.UserResponse {
	var tokoResponse *domain.TokoResponse
	if user.Toko.ID != 0 {
		tokoResponse = ToTokoResponse(&user.Toko)
	}

	return &domain.UserResponse{
		ID:                     user.ID,
		FName:                  user.FName,
		LName:                  user.LName,
		Email:                  user.Email,
		Alamat:                 user.Alamat,
		Thumbnail:              user.Thumbnail,
		Number:                 user.Number,
		Role:                   string(user.Role),
		Toko:                   tokoResponse,
		SubscriptionStatus:     user.SubscriptionStatus,
		SubscriptionExpiry:     FormatDate(user.SubscriptionExpiry),
		CreatedAt:              user.CreatedAt,
		SubscriptionExpiryDate: user.SubscriptionExpiry,
	}
}

func ToUserResponses(users []*domain.User) []*domain.UserResponse {
	var responses []*domain.UserResponse
	for _, user := range users {
		responses = append(responses, ToUserResponse(user))
	}
	return responses
}

func ToTokoResponse(toko *domain.Toko) *domain.TokoResponse {
	return &domain.TokoResponse{
		ID:      toko.ID,
		Name:    toko.Name,
		Address: toko.Address,
	}
}

func ToTokoResponses(tokos []*domain.Toko) []*domain.TokoResponse {

	var responses []*domain.TokoResponse

	for _, toko := range tokos {

		responses = append(responses, ToTokoResponse(toko))

	}

	return responses

}



// CheckoutRequest defines the structure for a checkout payload from the frontend

type CheckoutRequest struct {

	TokoID        uint       `json:"toko_id"`

	UserID        uint       `json:"user_id"`

	PaymentMethod string     `json:"payment_method"`

	Items         []CartItem `json:"items"`

}



// CartItem defines a single item in the shopping cart

type CartItem struct {

	ProductID uint `json:"product_id"`

	Quantity  int  `json:"quantity"`

}
