package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/exception"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ProductServiceImpl struct {
	ProductRepository         repository.ProductRepository
	ProductCategoryRepository repository.ProductCategoryRepository
	Validate                  *validator.Validate
}

func NewProductService(productRepository repository.ProductRepository, productCategoryRepository repository.ProductCategoryRepository,
	validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository:         productRepository,
		ProductCategoryRepository: productCategoryRepository,
		Validate:                  validate,
	}
}

func (p *ProductServiceImpl) Create(ctx context.Context, req *domain.ProductCreateRequest, file multipart.File, handler *multipart.FileHeader) error {
	if err := p.Validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
		return exception.BadRequest("validation failed: " + strings.Join(validationErrors, ", "))
	}
	EXP, err := helper.ParseDate(req.Exp)
	if err != nil {
		return exception.BadRequest("invalid date format, use yyyy-mm-dd or dd-mm-yyyy")
	}

	var thumbnailPath string
	if file != nil {
		if err := os.MkdirAll("static/image", os.ModePerm); err != nil {
			return exception.InternalServerError("failed to create image directory")
		}
		filePath := filepath.Join("static/image", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			return exception.InternalServerError("failed to save image")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			return exception.InternalServerError("failed to copy image file")
		}
		thumbnailPath = filePath
	}

	product := domain.Product{
		Name:       req.Name,
		Thumbnail:  thumbnailPath,
		Price:      req.Price,
		Exp:        EXP,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Barcode:    req.Barcode,
	}
	if _, err := p.ProductRepository.Create(ctx, &product); err != nil {
		return exception.InternalServerError("failed to create product")
	}
	return nil
}
func (p *ProductServiceImpl) FindAll(ctx context.Context) ([]*domain.ProductResponse, error) {
	product, category, err := p.ProductRepository.FindAll(ctx)
	if err != nil {
		return nil, exception.InternalServerError("data is not found")
	}
	if len(product) == 0 {
		return []*domain.ProductResponse{}, nil
	}
	return helper.ToProductResponses(product, category), nil
}
func (p *ProductServiceImpl) FindById(ctx context.Context, produkId int) (*domain.ProductResponse, error) {
	product, category, err := p.ProductRepository.FindById(ctx, uint(produkId))
	if err != nil {
		return nil, exception.NotFound("id category not exists")
	}
	return helper.ToProductResponse(product, category), nil
}
func (p *ProductServiceImpl) FindByBarcode(ctx context.Context, barcode string) (*domain.ProductResponse, error) {
	product, err := p.ProductRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return nil, exception.NotFound("barcode not exists")
	}
	return helper.ToProductResponse(product, &product.Category), nil
}
func (p *ProductServiceImpl) Update(ctx context.Context, req *domain.ProductUpdateRequest, file multipart.File, handler *multipart.FileHeader) error {
	if err := p.Validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
		return exception.BadRequest("validation failed: " + strings.Join(validationErrors, ", "))
	}

	product, _, err := p.ProductRepository.FindById(ctx, uint(req.ID))
	if err != nil {
		return exception.NotFound("id category not exists")
	}

	EXP, err := helper.ParseDate(req.Exp)
	if err != nil {
		return exception.BadRequest("invalid date format, use yyyy-mm-dd or dd-mm-yyyy")
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Exp = EXP
	product.Stock = req.Stock
	product.Barcode = req.Barcode
	product.CategoryID = req.CategoryID

	if file != nil {
		if err := os.MkdirAll("static/image", os.ModePerm); err != nil {
			return exception.InternalServerError("failed to create image directory")
		}
		filePath := filepath.Join("static/image", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			return exception.InternalServerError("failed to save image")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			return exception.InternalServerError("failed to copy image file")
		}
		product.Thumbnail = filePath
	}

	if _, err := p.ProductRepository.Update(ctx, product); err != nil {
		return exception.InternalServerError("failed to update product")
	}
	return nil
}
func (p *ProductServiceImpl) Delete(ctx context.Context, produkId int) error {
	product, _, err := p.ProductRepository.FindById(ctx, uint(produkId))
	if err != nil {
		return exception.NotFound("id product not exist")
	}
	if err := p.ProductRepository.Delete(ctx, product); err != nil {
		return exception.InternalServerError("failed to delete product")
	}
	return nil
}
func (p *ProductServiceImpl) UploadThumbnail(ctx context.Context, productId uint, file multipart.File, handler *multipart.FileHeader) error {
	defer file.Close()
	product, _, err := p.ProductRepository.FindById(ctx, productId)
	if err != nil {
		return exception.NotFound("id product not exists")
	}
	if err := os.MkdirAll("static/image", os.ModePerm); err != nil {
		return exception.InternalServerError("failed to create image directory")
	}
	filePath := filepath.Join("static/image", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return exception.InternalServerError("failed to save image")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return exception.InternalServerError("failed to copy image file")
	}
	if err := p.ProductRepository.UploadThumbnail(ctx, product.ID, filePath); err != nil {
		return exception.InternalServerError("failed to update product thumbnail")
	}
	return nil
}

func (p *ProductServiceImpl) UpdateStock(ctx context.Context, productId uint, req *domain.ProductUpdateStockRequest) error {
	if err := p.Validate.Struct(req); err != nil {
		return exception.BadRequest("invalid request, stock is required and must be zero or greater")
	}

	// Check if product exists
	if _, _, err := p.ProductRepository.FindById(ctx, productId); err != nil {
		return exception.NotFound("product not found")
	}

	if err := p.ProductRepository.UpdateStock(ctx, productId, req.Stock); err != nil {
		return exception.InternalServerError("failed to update stock")
	}

	return nil
}

func (p *ProductServiceImpl) FindLowStock(ctx context.Context, threshold uint) ([]*domain.ProductResponse, error) {
	products, err := p.ProductRepository.FindLowStock(ctx, threshold)
	if err != nil {
		return nil, exception.InternalServerError("failed to get low stock products")
	}

	var productResponses []*domain.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, helper.ToProductResponse(product, &product.Category))
	}

	return productResponses, nil
}
