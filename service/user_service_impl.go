package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

func (u *UserServiceImpl) Login(req *domain.LoginRequest) (string, error) {
	user, err := u.UserRepository.FindByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}
	token, err := helper.GenerateJWT(user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
