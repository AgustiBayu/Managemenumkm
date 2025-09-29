package helper

import "Managemenumkm/domain"

func ToProductResponse(product *domain.Product, category *domain.ProductCategory) *domain.ProductResponse {
	return &domain.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Thumbnail:  product.Thumbnail,
		Price:      product.Price,
		Exp:        FormatDate(product.Exp),
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
		ProductCategory: domain.ProductCategoryResponse{
			ID:       category.ID,
			Category: category.Category,
		},
		Barcode: product.Barcode,
	}
}
func ToProductResponses(producs []*domain.Product, categoriesMap map[uint]*domain.ProductCategory) []*domain.ProductResponse {
	var productResponses []*domain.ProductResponse
	for _, product := range producs {
		category, exits := categoriesMap[product.CategoryID]
		if !exits {
			category = &domain.ProductCategory{}
		}
		productResponses = append(productResponses, ToProductResponse(product, category))
	}
	return productResponses
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
