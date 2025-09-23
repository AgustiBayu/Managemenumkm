package helper

import "Managemenumkm/domain"

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
