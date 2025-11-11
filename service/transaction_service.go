package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
)

type TransactionService interface {
	Checkout(request helper.CheckoutRequest) (*domain.Transaction, error)
}
