package service

import (
	"Managemenumkm/domain"
	"context"
	"mime/multipart"
)

type UserService interface {
	Login(req *domain.LoginRequest) (string, error)
	Create(ctx context.Context, req *domain.UserCreateRequest) error
	FindAll(ctx context.Context, page, pageSize int) ([]*domain.UserResponse, int64, error)
	FindById(ctx context.Context, userID int) (*domain.UserResponse, error)
	FindByEmail(ctx context.Context, email string) (*domain.UserResponse, error)
	Update(ctx context.Context, req *domain.UserUpdateRequest, file multipart.File, handler *multipart.FileHeader) error
	Delete(ctx context.Context, userID int) error
	UploadThumbnail(ctx context.Context, userID uint, file multipart.File, handler *multipart.FileHeader) error
}
