package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
)

type TransactionService interface {
	Checkout(request helper.CheckoutRequest) (*domain.Transaction, error)
	CheckoutWithMember(request helper.CheckoutRequest, memberID uint) (*domain.Transaction, error)
}

type MemberCheckoutRequest struct {
	helper.CheckoutRequest
	MemberID uint `json:"member_id"`
}
