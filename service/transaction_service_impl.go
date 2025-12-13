package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionServiceImpl struct {
	db                    *gorm.DB
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
	memberService         MemberService
}

func NewTransactionService(db *gorm.DB, transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, memberService MemberService) TransactionService {
	return &TransactionServiceImpl{
		db:                    db,
		transactionRepository: transactionRepo,
		productRepository:     productRepo,
		memberService:         memberService,
	}
}

func (s *TransactionServiceImpl) Checkout(request helper.CheckoutRequest) (*domain.Transaction, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	var totalAmount uint
	var transactionItems []domain.TransactionItem

	for _, item := range request.Items {
		product, err := s.productRepository.FindById(item.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		var totalStock int
		for _, batch := range product.Batches {
			totalStock += int(batch.Stock)
		}

		if totalStock < item.Quantity {
			tx.Rollback()
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		subTotal := product.Price * uint(item.Quantity)
		totalAmount += subTotal

		transactionItems = append(transactionItems, domain.TransactionItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     float64(product.Price), // Price at the time of transaction
		})

		// Update stock
		remainingQty := item.Quantity
		for _, batch := range product.Batches {
			if remainingQty == 0 {
				break
			}
			if batch.Stock > 0 {
				deductQty := remainingQty
				if int(batch.Stock) < remainingQty {
					deductQty = int(batch.Stock)
				}
				newStock := int(batch.Stock) - deductQty
				err := s.productRepository.UpdateBatchStock(tx, batch.ID, uint(newStock))
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				remainingQty -= deductQty
			}
		}
	}

	transaction := &domain.Transaction{
		TokoID:          request.TokoID,
		UserID:          request.UserID,
		TransactionDate: time.Now(),
		TotalAmount:     float64(totalAmount),
		PaymentMethod:   request.PaymentMethod,
		Items:           transactionItems,
	}

	savedTransaction, err := s.transactionRepository.Save(tx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return savedTransaction, tx.Commit().Error
}

func (s *TransactionServiceImpl) CheckoutWithMember(request helper.CheckoutRequest, memberID uint) (*domain.Transaction, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Validate member exists
	_, err := s.memberService.GetMemberByID(memberID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("member not found")
	}

	var totalAmount uint
	var transactionItems []domain.TransactionItem

	for _, item := range request.Items {
		product, err := s.productRepository.FindById(item.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		var totalStock int
		for _, batch := range product.Batches {
			totalStock += int(batch.Stock)
		}

		if totalStock < item.Quantity {
			tx.Rollback()
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		subTotal := product.Price * uint(item.Quantity)
		totalAmount += subTotal

		transactionItems = append(transactionItems, domain.TransactionItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     float64(product.Price), // Price at the time of transaction
		})

		// Update stock
		remainingQty := item.Quantity
		for _, batch := range product.Batches {
			if remainingQty == 0 {
				break
			}
			if batch.Stock > 0 {
				deductQty := remainingQty
				if int(batch.Stock) < remainingQty {
					deductQty = int(batch.Stock)
				}
				newStock := int(batch.Stock) - deductQty
				err := s.productRepository.UpdateBatchStock(tx, batch.ID, uint(newStock))
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				remainingQty -= deductQty
			}
		}
	}

	// No member tier discount - all members get the same base rate
	var finalAmount = float64(totalAmount)
	// Discount logic moved to separate discount system

	transaction := &domain.Transaction{
		TokoID:          request.TokoID,
		UserID:          request.UserID,
		TransactionDate: time.Now(),
		TotalAmount:     finalAmount,
		PaymentMethod:   request.PaymentMethod,
		Items:           transactionItems,
	}

	savedTransaction, err := s.transactionRepository.Save(tx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Award points to member after successful transaction
	if finalAmount > 0 {
		_, err = s.memberService.EarnPoints(
			memberID,
			savedTransaction.ID,
			finalAmount,
			fmt.Sprintf("Purchase transaction #%d", savedTransaction.ID),
		)
		if err != nil {
			// Log error but don't fail the transaction
			// In a production environment, you might want to retry this or handle it differently
			// For now, we'll just log it and continue
			// TODO: Implement proper logging and retry mechanism
			_ = err
		}
	}

	return savedTransaction, tx.Commit().Error
}
