package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/exception"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type productServiceImpl struct {
	ProductRepo      repository.ProductRepository
	ProductBatchRepo repository.ProductBatchRepository
	CategoryRepo     repository.ProductCategoryRepository // Keep for validation
	Validate         *validator.Validate
}

func NewProductService(productRepo repository.ProductRepository, batchRepo repository.ProductBatchRepository, categoryRepo repository.ProductCategoryRepository, validate *validator.Validate) ProductService {
	return &productServiceImpl{
		ProductRepo:      productRepo,
		ProductBatchRepo: batchRepo,
		CategoryRepo:     categoryRepo,
		Validate:         validate,
	}
}

// Create handles creation of a new master product.
func (s *productServiceImpl) Create(req domain.ProductCreateRequest) (domain.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.ProductResponse{}, exception.BadRequest(err.Error())
	}

	// Check if SKU is unique
	_, err := s.ProductRepo.FindBySKU(req.SKU)
	if err == nil {
		return domain.ProductResponse{}, exception.InternalServerError("Product with this SKU already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ProductResponse{}, err // Handle other db errors
	}

	// Check if category exists
	ctx := context.Background()
	_, err = s.CategoryRepo.FindById(ctx, int(req.CategoryID))
	if err != nil {
		return domain.ProductResponse{}, exception.NotFound("Product category not found")
	}

	product := domain.Product{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}

	newProduct, err := s.ProductRepo.Save(product)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	// Reload to get associations
	newProduct, _ = s.ProductRepo.FindById(newProduct.ID)
	return helper.ToProductResponse(newProduct), nil
}

// Update handles updates to a master product.
func (s *productServiceImpl) Update(req domain.ProductUpdateRequest) (domain.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.ProductResponse{}, exception.BadRequest(err.Error())
	}

	// Check if product exists
	product, err := s.ProductRepo.FindById(req.ID)
	if err != nil {
		return domain.ProductResponse{}, exception.NotFound("Product not found")
	}

	// Check if SKU is being changed to one that already exists
	if product.SKU != req.SKU {
		_, err := s.ProductRepo.FindBySKU(req.SKU)
		if err == nil {
			return domain.ProductResponse{}, exception.InternalServerError("Another product with this SKU already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ProductResponse{}, err
		}
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.Price = req.Price
	product.CategoryID = req.CategoryID

	updatedProduct, err := s.ProductRepo.Update(product)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	// Reload to get associations
	updatedProduct, _ = s.ProductRepo.FindById(updatedProduct.ID)
	return helper.ToProductResponse(updatedProduct), nil
}

// AddStock creates a new product batch.
func (s *productServiceImpl) AddStock(req domain.AddStockStep2Request) (domain.ProductBatchResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.ProductBatchResponse{}, exception.BadRequest(err.Error())
	}

	// Check if barcode is already used
	_, err := s.ProductBatchRepo.FindByBarcode(req.Barcode)
	if err == nil {
		return domain.ProductBatchResponse{}, exception.InternalServerError("This barcode is already registered in another batch")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ProductBatchResponse{}, err
	}

	// Check if master product exists
	_, err = s.ProductRepo.FindById(req.ProductID)
	if err != nil {
		return domain.ProductBatchResponse{}, exception.NotFound("Master product not found")
	}

	expDate, err := helper.ParseDate(req.Exp)
	if err != nil {
		return domain.ProductBatchResponse{}, exception.BadRequest("Invalid expiry date format. Use YYYY-MM-DD")
	}

	batch := domain.ProductBatch{
		ProductID: req.ProductID,
		Barcode:   req.Barcode,
		Exp:       expDate,
		Stock:     req.Stock,
		DateAdded: time.Now(),
	}

	newBatch, err := s.ProductBatchRepo.Save(batch)
	if err != nil {
		return domain.ProductBatchResponse{}, err
	}

	return helper.ToProductBatchResponse(newBatch), nil
}

func (s *productServiceImpl) Delete(productID uint) error {
	// Optional: Check if product has batches and prevent deletion if it does
	return s.ProductRepo.Delete(productID)
}

func (s *productServiceImpl) FindById(productID uint) (domain.ProductResponse, error) {
	product, err := s.ProductRepo.FindById(productID)
	if err != nil {
		return domain.ProductResponse{}, exception.NotFound("Product not found")
	}
	return helper.ToProductResponse(product), nil
}

func (s *productServiceImpl) FindAll() ([]domain.ProductResponse, error) {
	products, err := s.ProductRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return helper.ToProductResponses(products), nil
}
