package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/exception"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"context"

	"github.com/go-playground/validator/v10"
)

type ProductCategoryServiceImpl struct {
	ProductCategoryRepository repository.ProductCategoryRepository
	Validate                  *validator.Validate
}

func NewProductCategoryService(productCategoryRepository repository.ProductCategoryRepository,
	validate *validator.Validate) ProductCategoryService {
	return &ProductCategoryServiceImpl{
		ProductCategoryRepository: productCategoryRepository,
		Validate:                  validate,
	}
}

func (p *ProductCategoryServiceImpl) Create(ctx context.Context, req *domain.ProductCategoryCreateRequest) error {
	if err := p.Validate.Struct(req); err != nil {
		return exception.BadRequest("field is not valid")
	}
	category := domain.ProductCategory{
		Category: req.Category,
	}
	if _, err := p.ProductCategoryRepository.Create(ctx, &category); err != nil {
		return exception.InternalServerError("failed to create category")
	}
	return nil
}
func (p *ProductCategoryServiceImpl) FindAll(ctx context.Context) ([]*domain.ProductCategoryResponse, error) {
	categories, err := p.ProductCategoryRepository.FindAll(ctx)
	if err != nil {
		return nil, exception.InternalServerError("data is not exist")
	}
	if len(categories) == 0 {
		return []*domain.ProductCategoryResponse{}, nil
	}
	return helper.ToProductCategoryResponses(categories), nil
}
func (p *ProductCategoryServiceImpl) FindById(ctx context.Context, categoryID int) (*domain.ProductCategoryResponse, error) {
	category, err := p.ProductCategoryRepository.FindById(ctx, categoryID)
	if err != nil {
		return nil, exception.NotFound("id category not exist")
	}
	return helper.ToProductCategoryResponse(category), nil
}
func (p *ProductCategoryServiceImpl) Update(ctx context.Context, req *domain.ProductCategoryUpdateRequest) error {
	if err := p.Validate.Struct(req); err != nil {
		return exception.BadRequest("field is not valid")
	}
	category, err := p.ProductCategoryRepository.FindById(ctx, int(req.ID))
	if err != nil {
		return exception.NotFound("id category not exist")
	}
	category.Category = req.Category
	if _, err := p.ProductCategoryRepository.Update(ctx, category); err != nil {
		return exception.InternalServerError("failed to update category")
	}
	return nil
}
func (p *ProductCategoryServiceImpl) Delete(ctx context.Context, categoryID int) error {
	category, err := p.ProductCategoryRepository.FindById(ctx, categoryID)
	if err != nil {
		return exception.NotFound("id category not exist")
	}
	if err := p.ProductCategoryRepository.Delete(ctx, category); err != nil {
		return exception.InternalServerError("failed to delete category")
	}
	return nil
}
