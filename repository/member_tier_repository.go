package repository

import (
	"Managemenumkm/domain"
	"gorm.io/gorm"
)

type MemberTierRepository interface {
	Create(tier domain.MemberTier) (domain.MemberTier, error)
	Update(tier domain.MemberTier) (domain.MemberTier, error)
	Delete(tierID uint) error
	FindByID(tierID uint) (domain.MemberTier, error)
	FindByTokoID(tokoID uint) ([]domain.MemberTier, error)
	FindByPoints(points int, tokoID uint) (domain.MemberTier, error)
	FindBySpentAmount(amount float64, tokoID uint) (domain.MemberTier, error)
	GetDefaultTier(tokoID uint) (domain.MemberTier, error)
}

type MemberTierRepositoryImpl struct {
	DB *gorm.DB
}

func NewMemberTierRepository(db *gorm.DB) MemberTierRepository {
	return &MemberTierRepositoryImpl{DB: db}
}

func (r *MemberTierRepositoryImpl) Create(tier domain.MemberTier) (domain.MemberTier, error) {
	err := r.DB.Create(&tier).Error
	return tier, err
}

func (r *MemberTierRepositoryImpl) Update(tier domain.MemberTier) (domain.MemberTier, error) {
	err := r.DB.Save(&tier).Error
	return tier, err
}

func (r *MemberTierRepositoryImpl) Delete(tierID uint) error {
	err := r.DB.Delete(&domain.MemberTier{}, tierID).Error
	return err
}

func (r *MemberTierRepositoryImpl) FindByID(tierID uint) (domain.MemberTier, error) {
	var tier domain.MemberTier
	err := r.DB.First(&tier, tierID).Error
	return tier, err
}

func (r *MemberTierRepositoryImpl) FindByTokoID(tokoID uint) ([]domain.MemberTier, error) {
	var tiers []domain.MemberTier
	err := r.DB.Where("toko_id = ?", tokoID).Order("min_spent ASC, min_points ASC").Find(&tiers).Error
	return tiers, err
}

func (r *MemberTierRepositoryImpl) FindByPoints(points int, tokoID uint) (domain.MemberTier, error) {
	var tier domain.MemberTier
	err := r.DB.Where("toko_id = ? AND min_points <= ? AND is_active = ?", tokoID, points, true).
		Order("min_points DESC").First(&tier).Error
	return tier, err
}

func (r *MemberTierRepositoryImpl) FindBySpentAmount(amount float64, tokoID uint) (domain.MemberTier, error) {
	var tier domain.MemberTier
	err := r.DB.Where("toko_id = ? AND min_spent <= ? AND is_active = ?", tokoID, amount, true).
		Order("min_spent DESC").First(&tier).Error
	return tier, err
}

func (r *MemberTierRepositoryImpl) GetDefaultTier(tokoID uint) (domain.MemberTier, error) {
	var tier domain.MemberTier
	err := r.DB.Where("toko_id = ? AND is_default = ? AND is_active = ?", tokoID, true, true).First(&tier).Error
	return tier, err
}