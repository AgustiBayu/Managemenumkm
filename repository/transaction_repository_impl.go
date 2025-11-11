package repository

import (
	"Managemenumkm/domain"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

// Save saves a transaction and its items within a given database transaction
func (r *TransactionRepositoryImpl) Save(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error) {
	err := tx.Create(&transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
