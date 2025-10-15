package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type productBatchServiceImpl struct {
	ProductBatchRepository repository.ProductBatchRepository
	Validate               *validator.Validate
}

func NewProductBatchService(batchRepo repository.ProductBatchRepository, validate *validator.Validate) ProductBatchService {
	return &productBatchServiceImpl{ProductBatchRepository: batchRepo, Validate: validate}
}

func (s *productBatchServiceImpl) FindByBarcode(barcode string) (domain.ProductBatchResponse, error) {
	batch, err := s.ProductBatchRepository.FindByBarcode(barcode)
	if err != nil {
		return domain.ProductBatchResponse{}, err // Consider custom errors
	}

	return helper.ToProductBatchResponse(batch), nil
}

func (s *productBatchServiceImpl) FindByProductID(productID uint) ([]domain.ProductBatchResponse, error) {
	batches, err := s.ProductBatchRepository.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	return helper.ToProductBatchResponses(batches), nil
}

func (s *productBatchServiceImpl) FindById(batchID uint) (domain.ProductBatchResponse, error) {
	batch, err := s.ProductBatchRepository.FindById(batchID)
	if err != nil {
		return domain.ProductBatchResponse{}, err
	}

	return helper.ToProductBatchResponse(batch), nil
}

func (s *productBatchServiceImpl) Update(request domain.ProductBatchUpdateRequest) (domain.ProductBatchResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return domain.ProductBatchResponse{}, err
	}

	// Get existing batch to ensure we don't change other fields like barcode
	existingBatch, err := s.ProductBatchRepository.FindById(request.ID)
	if err != nil {
		return domain.ProductBatchResponse{}, err
	}

	// Update fields from request
	expTime, _ := time.Parse("2006-01-02", request.Exp) // Error handled by validation
	existingBatch.Stock = request.Stock
	existingBatch.Exp = expTime

	updatedBatch, err := s.ProductBatchRepository.Update(existingBatch)
	if err != nil {
		return domain.ProductBatchResponse{}, err
	}

	return helper.ToProductBatchResponse(updatedBatch), nil
}
