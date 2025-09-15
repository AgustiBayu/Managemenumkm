package repository

import (
	"Managemenumkm/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error)
	FindByEmail(email string) (domain.User, error)
	FindByIdAndFNameAndLNameAndEmailAndAlamat(ctx context.Context, userID uint, fname, lname, email, alamat string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, userID uint) error
}
