package repository

import (
	"Managemenumkm/domain"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Save(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error)
}
