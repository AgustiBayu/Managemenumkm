package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/repository"
	"errors"

	"github.com/go-playground/validator/v10"
)

type MemberTierServiceImpl struct {
	MemberTierRepository repository.MemberTierRepository
	Validate             *validator.Validate
}

func NewMemberTierService(memberTierRepository repository.MemberTierRepository, validate *validator.Validate) MemberTierService {
	return &MemberTierServiceImpl{
		MemberTierRepository: memberTierRepository,
		Validate:             validate,
	}
}

func (s *MemberTierServiceImpl) CreateMemberTier(request domain.MemberTier) (domain.MemberTier, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return domain.MemberTier{}, err
	}

	// Check if this is being set as default tier
	if request.IsDefault {
		// Unset any existing default tier for this toko
		existingDefault, err := s.MemberTierRepository.GetDefaultTier(request.TokoID)
		if err == nil && existingDefault.ID != 0 {
			existingDefault.IsDefault = false
			s.MemberTierRepository.Update(existingDefault)
		}
	}

	tier, err := s.MemberTierRepository.Create(request)
	if err != nil {
		return domain.MemberTier{}, err
	}

	return tier, nil
}

func (s *MemberTierServiceImpl) UpdateMemberTier(tierID uint, request domain.MemberTier) (domain.MemberTier, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return domain.MemberTier{}, err
	}

	existingTier, err := s.MemberTierRepository.FindByID(tierID)
	if err != nil {
		return domain.MemberTier{}, err
	}

	// Check if this is being set as default tier
	if request.IsDefault && !existingTier.IsDefault {
		// Unset any existing default tier for this toko
		existingDefault, err := s.MemberTierRepository.GetDefaultTier(request.TokoID)
		if err == nil && existingDefault.ID != tierID {
			existingDefault.IsDefault = false
			s.MemberTierRepository.Update(existingDefault)
		}
	}

	existingTier.Name = request.Name
	existingTier.Description = request.Description
	existingTier.MinPoints = request.MinPoints
	existingTier.MinSpent = request.MinSpent
	existingTier.DiscountRate = request.DiscountRate
	existingTier.PointRate = request.PointRate
	existingTier.IsDefault = request.IsDefault
	existingTier.IsActive = request.IsActive

	updatedTier, err := s.MemberTierRepository.Update(existingTier)
	if err != nil {
		return domain.MemberTier{}, err
	}

	return updatedTier, nil
}

func (s *MemberTierServiceImpl) DeleteMemberTier(tierID uint) error {
	err := s.MemberTierRepository.Delete(tierID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MemberTierServiceImpl) GetMemberTierByID(tierID uint) (domain.MemberTier, error) {
	tier, err := s.MemberTierRepository.FindByID(tierID)
	if err != nil {
		return domain.MemberTier{}, err
	}
	return tier, nil
}

func (s *MemberTierServiceImpl) GetMemberTiersByTokoID(tokoID uint) ([]domain.MemberTier, error) {
	tiers, err := s.MemberTierRepository.FindByTokoID(tokoID)
	if err != nil {
		return nil, err
	}
	return tiers, nil
}

func (s *MemberTierServiceImpl) InitializeDefaultTiers(tokoID uint) ([]domain.MemberTier, error) {
	// Check if tiers already exist for this toko
	existingTiers, err := s.MemberTierRepository.FindByTokoID(tokoID)
	if err == nil && len(existingTiers) > 0 {
		return existingTiers, nil // Tiers already exist
	}

	// Create default tiers
	defaultTiers := []domain.MemberTier{
		{
			TokoID:       tokoID,
			Name:         "Bronze",
			Description:  "Basic member tier",
			MinPoints:    0,
			MinSpent:     0,
			DiscountRate: 0.00, // 0% discount
			PointRate:    1.0,  // 1 point per Rp 100
			IsDefault:    true,
			IsActive:     true,
		},
		{
			TokoID:       tokoID,
			Name:         "Silver",
			Description:  "Silver member tier",
			MinPoints:    1000,
			MinSpent:     500000, // Rp 500,000
			DiscountRate: 0.05,   // 5% discount
			PointRate:    1.2,    // 1.2 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
		{
			TokoID:       tokoID,
			Name:         "Gold",
			Description:  "Gold member tier",
			MinPoints:    5000,
			MinSpent:     2000000, // Rp 2,000,000
			DiscountRate: 0.10,    // 10% discount
			PointRate:    1.5,     // 1.5 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
		{
			TokoID:       tokoID,
			Name:         "Platinum",
			Description:  "Platinum member tier",
			MinPoints:    10000,
			MinSpent:     5000000, // Rp 5,000,000
			DiscountRate: 0.15,    // 15% discount
			PointRate:    2.0,     // 2 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
	}

	var createdTiers []domain.MemberTier

	for _, tier := range defaultTiers {
		createdTier, err := s.MemberTierRepository.Create(tier)
		if err != nil {
			return createdTiers, errors.New("failed to create default tier: " + tier.Name)
		}
		createdTiers = append(createdTiers, createdTier)
	}

	return createdTiers, nil
}