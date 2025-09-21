package repository

import (
	"Managemenumkm/domain"
	"context"
)

type TokoRepository interface {
	Create(ctx context.Context, toko *domain.Toko) (*domain.Toko, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*domain.Toko, int64, error)
	FindById(ctx context.Context, tokoID int) (*domain.Toko, error)
	Update(ctx context.Context, toko *domain.Toko) (*domain.Toko, error)
	Delete(ctx context.Context, toko *domain.Toko) error
}
