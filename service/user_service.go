package service

import (
	"Managemenumkm/domain"
)

type UserService interface {
	Login(req *domain.LoginRequest) (string, error)
}
