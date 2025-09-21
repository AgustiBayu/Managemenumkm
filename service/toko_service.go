package service

import (
	"Managemenumkm/domain"
	"context"
)

type TokoService interface {
	Create(ctx context.Context, req *domain.TokoCreateRequest) error
	FindAll(ctx context.Context, page, pageSize int) ([]*domain.TokoResponse, int64, error)
	FindById(ctx context.Context, tokoID int) (*domain.TokoResponse, error)
	Update(ctx context.Context, req *domain.TokoUpdateRequest) error
	Delete(ctx context.Context, tokoID int) error
}
