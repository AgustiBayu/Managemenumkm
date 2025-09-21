package repository

import (
	"Managemenumkm/domain"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (u *UserRepositoryImpl) FindById(ctx context.Context, userID int) (*domain.User, error) {
	var user *domain.User
	if err := u.DB.WithContext(ctx).Preload("Toko").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user *domain.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func (u *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("create user failed err: %s", err)
	}
	return user, nil
}
func (u *UserRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	var user []*domain.User
	var totalItem int64

	if err := u.DB.WithContext(ctx).Model(&domain.User{}).Count(&totalItem).Error; err != nil {
		return nil, 0, fmt.Errorf("user is null: %s", err)
	}
	offset := (page - 1) * pageSize
	if err := u.DB.WithContext(ctx).Preload("Toko").Limit(pageSize).Offset(offset).Find(&user).Error; err != nil {
		return nil, 0, err
	}
	return user, totalItem, nil
}
func (u *UserRepositoryImpl) FindByIdAndFNameAndLNameAndEmailAndAlamat(ctx context.Context, userID uint, fname, lname, email, alamat string) (*domain.User, error) {
	var user *domain.User
	hasil := u.DB.WithContext(ctx).Where("id = ? AND first_name = ? AND last_name = ? AND email = ? AND alamat = ?", userID, fname, lname, email, alamat).First(&user)
	if hasil != nil {
		if errors.Is(hasil.Error, gorm.ErrRecordNotFound) {
			return nil, hasil.Error
		}
		return nil, hasil.Error
	}
	return user, nil
}
func (u *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := u.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed update user err: %s", err)
	}
	return user, nil
}
func (u *UserRepositoryImpl) Delete(ctx context.Context, user *domain.User) error {
	if err := u.DB.WithContext(ctx).Delete(&user).Error; err != nil {
		return fmt.Errorf("failed delete user err: %s", err)
	}
	return nil
}

func (p *UserRepositoryImpl) UploadThumbnail(ctx context.Context, userID uint, path string) error {
	var user *domain.User
	if err := p.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	user.Thumbnail = path
	if err := p.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryImpl) FindAllByToko(ctx context.Context, tokoID uint, page, pageSize int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var totalItem int64

	// Menghitung total item untuk toko tertentu
	if err := u.DB.WithContext(ctx).Model(&domain.User{}).Where("toko_id = ?", tokoID).Count(&totalItem).Error; err != nil {
		return nil, 0, fmt.Errorf("users not found in toko: %w", err)
	}

	// Mengambil data pengguna dengan paginasi
	offset := (page - 1) * pageSize
	if err := u.DB.WithContext(ctx).Where("toko_id = ?", tokoID).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItem, nil
}
