package service

import (
	"Managemenumkm/domain"
)

type MemberTierService interface {
	CreateMemberTier(request domain.MemberTier) (domain.MemberTier, error)
	UpdateMemberTier(tierID uint, request domain.MemberTier) (domain.MemberTier, error)
	DeleteMemberTier(tierID uint) error
	GetMemberTierByID(tierID uint) (domain.MemberTier, error)
	GetMemberTiersByTokoID(tokoID uint) ([]domain.MemberTier, error)
	InitializeDefaultTiers(tokoID uint) ([]domain.MemberTier, error)
}

type MemberTierServiceRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	MinPoints    int     `json:"min_points"`
	MinSpent     float64 `json:"min_spent"`
	DiscountRate float64 `json:"discount_rate"`
	PointRate    float64 `json:"point_rate"`
	IsDefault    bool    `json:"is_default"`
}