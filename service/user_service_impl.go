package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/exception"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (u *UserServiceImpl) Login(req *domain.LoginRequest) (string, error) {
	user, err := u.UserRepository.FindByEmail(req.Email)
	if err != nil {
		return "", exception.InternalServerError("invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", exception.InternalServerError("invalid email or password")
	}
	token, err := helper.GenerateJWT(user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserServiceImpl) Create(ctx context.Context, req *domain.UserCreateRequest) error {
	if err := u.Validate.Struct(req); err != nil {
		return exception.BadRequest("field not valid")
	}
	EXP, err := helper.ParseDate(req.SubscriptionExpiry)
	if err != nil {
		return exception.BadRequest("invalid date format, use yyyy-mm-dd or dd-mm-yyyy")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return exception.InternalServerError("failed to hash password")
	}

	user := domain.User{
		FName:              req.FName,
		LName:              req.LName,
		Email:              req.Email,
		Password:           string(hashedPassword), // Save hashed password
		Alamat:             req.Alamat,
		Thumbnail:          "",
		Number:             req.Number,
		Role:               req.Role,
		TokoID:             req.TokoID,
		SubscriptionStatus: req.SubscriptionStatus,
		SubscriptionExpiry: EXP,
	}
	if _, err := u.UserRepository.Create(ctx, &user); err != nil {
		return exception.InternalServerError("failed create user")
	}
	return nil
}
func (u *UserServiceImpl) FindAll(ctx context.Context, page, pageSize int) ([]*domain.UserResponse, int64, error) {
	users, totalItems, err := u.UserRepository.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, exception.InternalServerError("read user not valid")
	}
	return helper.ToUserResponses(users), totalItems, nil
}
func (u *UserServiceImpl) FindById(ctx context.Context, userID int) (*domain.UserResponse, error) {
	user, err := u.UserRepository.FindById(ctx, userID)
	if err != nil {
		return nil, exception.NotFound("id user not exists")
	}
	return helper.ToUserResponse(user), nil
}

func (u *UserServiceImpl) FindByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	user, err := u.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, exception.NotFound("email user not exists")
	}
	return helper.ToUserResponse(user), nil
}
func (u *UserServiceImpl) Update(ctx context.Context, req *domain.UserUpdateRequest, file multipart.File, handler *multipart.FileHeader) error {
	if err := u.Validate.Struct(req); err != nil {
		return exception.BadRequest("field not valid")
	}
	user, err := u.UserRepository.FindById(ctx, int(req.ID))
	if err != nil {
		return exception.NotFound("id user not exists")
	}
	EXP, err := helper.ParseDate(req.SubscriptionExpiry)
	if err != nil {
		return exception.BadRequest("invalid date format, use yyyy-mm-dd or dd-mm-yyyy")
	}
	user.FName = req.FName
	user.LName = req.LName
	user.Email = req.Email

	// Update password only if a new one is provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return exception.InternalServerError("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}

	user.Alamat = req.Alamat
	user.Number = req.Number
	user.Role = req.Role
	user.TokoID = req.TokoID
	user.SubscriptionStatus = req.SubscriptionStatus
	user.SubscriptionExpiry = EXP

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
		user.Thumbnail = filePath
	}

	if _, err := u.UserRepository.Update(ctx, user); err != nil {
		return exception.InternalServerError("failed to update user")
	}
	return nil
}
func (u *UserServiceImpl) Delete(ctx context.Context, userID int) error {
	user, err := u.UserRepository.FindById(ctx, userID)
	if err != nil {
		return exception.NotFound("id user not exists")
	}
	if err := u.UserRepository.Delete(ctx, user); err != nil {
		return exception.InternalServerError("failed to delete product")
	}
	return nil
}

func (u *UserServiceImpl) UploadThumbnail(ctx context.Context, userID uint, file multipart.File, handler *multipart.FileHeader) error {
	defer file.Close()
	user, err := u.UserRepository.FindById(ctx, int(userID))
	if err != nil {
		return exception.NotFound("id user not exists")
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
	if err := u.UserRepository.UploadThumbnail(ctx, user.ID, filePath); err != nil {
		return exception.InternalServerError("failed to update user thumbnail")
	}
	return nil
}
