package service

import (
	"Managemenumkm/domain"
	"context"
	"mime/multipart"
)

type ProductService interface {
	Create(ctx context.Context, req *domain.ProductCreateRequest, file multipart.File, handler *multipart.FileHeader) error
	FindAll(ctx context.Context) ([]*domain.ProductResponse, error)
	FindById(ctx context.Context, produkId int) (*domain.ProductResponse, error)
	FindByBarcode(ctx context.Context, barcode string) (*domain.ProductResponse, error)
	FindLowStock(ctx context.Context, threshold uint) ([]*domain.ProductResponse, error)
	Update(ctx context.Context, req *domain.ProductUpdateRequest, file multipart.File, handler *multipart.FileHeader) error
	UpdateStock(ctx context.Context, productId uint, req *domain.ProductUpdateStockRequest) error
	Delete(ctx context.Context, produkId int) error
}
