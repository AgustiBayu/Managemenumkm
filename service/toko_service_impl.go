package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/exception"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"context"

	"github.com/go-playground/validator/v10"
)

type TokoServiceImpl struct {
	TokoRepository repository.TokoRepository
	Validate       *validator.Validate
}

func NewTokoService(tokoRepository repository.TokoRepository, validate *validator.Validate) TokoService {
	return &TokoServiceImpl{
		TokoRepository: tokoRepository,
		Validate:       validate,
	}
}

func (s *TokoServiceImpl) Create(ctx context.Context, req *domain.TokoCreateRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return exception.BadRequest("field not valid")
	}

	toko := domain.Toko{
		Name:    req.Name,
		Address: req.Address,
	}

	if _, err := s.TokoRepository.Create(ctx, &toko); err != nil {
		return exception.InternalServerError("failed create toko")
	}

	return nil
}

func (s *TokoServiceImpl) FindAll(ctx context.Context, page, pageSize int) ([]*domain.TokoResponse, int64, error) {
	tokos, totalItems, err := s.TokoRepository.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, exception.InternalServerError("read toko not valid")
	}

	return helper.ToTokoResponses(tokos), totalItems, nil
}

func (s *TokoServiceImpl) FindById(ctx context.Context, tokoID int) (*domain.TokoResponse, error) {
	toko, err := s.TokoRepository.FindById(ctx, tokoID)
	if err != nil {
		return nil, exception.NotFound("id toko not exists")
	}

	return helper.ToTokoResponse(toko), nil
}

func (s *TokoServiceImpl) Update(ctx context.Context, req *domain.TokoUpdateRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return exception.BadRequest("field not valid")
	}

	toko, err := s.TokoRepository.FindById(ctx, int(req.ID))
	if err != nil {
		return exception.NotFound("id toko not exists")
	}

	toko.Name = req.Name
	toko.Address = req.Address

	if _, err := s.TokoRepository.Update(ctx, toko); err != nil {
		return exception.InternalServerError("failed to update toko")
	}

	return nil
}

func (s *TokoServiceImpl) Delete(ctx context.Context, tokoID int) error {
	toko, err := s.TokoRepository.FindById(ctx, tokoID)
	if err != nil {
		return exception.NotFound("id toko not exists")
	}

	if err := s.TokoRepository.Delete(ctx, toko); err != nil {
		return exception.InternalServerError("failed to delete toko")
	}

	return nil
}
