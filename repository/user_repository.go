package repository

import (
	"Managemenumkm/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error)
	FindAllByToko(ctx context.Context, tokoID uint, page, pageSize int) ([]*domain.User, int64, error)
	FindByEmail(email string) (*domain.User, error)
	FindById(ctx context.Context, userID int) (*domain.User, error)
	FindByIdAndFNameAndLNameAndEmailAndAlamat(ctx context.Context, userID uint, fname, lname, email, alamat string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, user *domain.User) error
	UploadThumbnail(ctx context.Context, userID uint, path string) error
	CountByTokoIDAndRole(ctx context.Context, tokoID uint, role domain.Role) (int64, error)
}
