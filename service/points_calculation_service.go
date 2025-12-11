package service

import (
	"context"
	"Managemenumkm/domain"
	"Managemenumkm/repository"
	"errors"
	"fmt"
	"time"
)

type PointsCalculationService interface {
	CalculatePoints(ctx context.Context, transaction *domain.Transaction, member *domain.Member) (int, error)
	CalculatePointsValue(ctx context.Context, points int, memberTier string) (float64, error)
	ProcessPointsRedemption(ctx context.Context, memberID uint, pointsToRedeem int) (*domain.MemberTransaction, error)
	GetPointsExpiringSoon(ctx context.Context, memberID uint, days int) ([]domain.MemberTransaction, error)
	ExpireOldPoints(ctx context.Context) error
}

type pointsCalculationServiceImpl struct {
	memberRepo   repository.MemberRepository
	pointRules   PointRules
}

type PointRules struct {
	BaseRate        float64  // Points per Rupiah (e.g., 0.0001 = 1 point per Rp 10,000)
	TierMultipliers map[string]float64  // Tier bonus multipliers
	BonusCategories map[string]float64  // Category-specific multipliers
	MinPointsRedeem int       // Minimum points for redemption
	PointsPerRupiah float64   // Conversion rate for points to money
}

func NewPointsCalculationService(memberRepo repository.MemberRepository) PointsCalculationService {
	pointRules := PointRules{
		BaseRate:        0.0001, // 1 point per Rp 10,000
		TierMultipliers: map[string]float64{
			"Bronze":    1.0,
			"Silver":    1.5,
			"Gold":      2.0,
			"Platinum":  3.0,
		},
		BonusCategories: map[string]float64{
			"ELECTRONICS": 2.0, // 2x points for electronics
			"FASHION":     1.5, // 1.5x points for fashion
		},
		MinPointsRedeem: 100,
		PointsPerRupiah: 100.0, // 100 points = Rp 1,000
	}

	return &pointsCalculationServiceImpl{
		memberRepo: memberRepo,
		pointRules: pointRules,
	}
}

func (s *pointsCalculationServiceImpl) CalculatePoints(ctx context.Context, transaction *domain.Transaction, member *domain.Member) (int, error) {
	if member == nil {
		return 0, nil
	}

	totalPoints := 0

	// Calculate base points for each item
	for _, item := range transaction.Items {
		itemSubtotal := float64(item.Quantity) * item.Price
		itemPoints := int(itemSubtotal * s.pointRules.BaseRate)

		// Apply tier multiplier
		tierMultiplier := s.pointRules.TierMultipliers[member.MemberTier.Name]
		itemPoints = int(float64(itemPoints) * tierMultiplier)

		// Apply category bonus if applicable (placeholder logic)
		// You would need to load product with category to implement this
		totalPoints += itemPoints
	}

	// Round down to nearest integer
	if totalPoints < 1 {
		totalPoints = 0
	}

	return totalPoints, nil
}

func (s *pointsCalculationServiceImpl) CalculatePointsValue(ctx context.Context, points int, memberTier string) (float64, error) {
	// Base conversion rate
	value := float64(points) / s.pointRules.PointsPerRupiah

	// Apply tier bonus
	if tierMultiplier, exists := s.pointRules.TierMultipliers[memberTier]; exists {
		value = value * tierMultiplier
	}

	return value, nil
}

func (s *pointsCalculationServiceImpl) ProcessPointsRedemption(ctx context.Context, memberID uint, pointsToRedeem int) (*domain.MemberTransaction, error) {
	if pointsToRedeem < s.pointRules.MinPointsRedeem {
		return nil, fmt.Errorf("minimum points to redeem is %d", s.pointRules.MinPointsRedeem)
	}

	// Get member current points
	member, err := s.memberRepo.FindByID(memberID)
	if err != nil {
		return nil, err
	}

	if member.TotalPoints < pointsToRedeem {
		return nil, errors.New("insufficient points balance")
	}

	// Create points transaction
	transaction := domain.MemberTransaction{
		MemberID:     memberID,
		Points:       -pointsToRedeem, // Negative for redemption
		Type:         "REDEEM",
		Description:  "Points redemption for discount",
		BalanceAfter: member.TotalPoints - pointsToRedeem,
		CreatedAt:    time.Now(),
	}

	// Update member balance
	newBalance := member.TotalPoints - pointsToRedeem
	err = s.memberRepo.UpdatePoints(memberID, newBalance)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *pointsCalculationServiceImpl) GetPointsExpiringSoon(ctx context.Context, memberID uint, days int) ([]domain.MemberTransaction, error) {
	// This would need additional implementation in the repository
	// For now, return empty slice
	return []domain.MemberTransaction{}, nil
}

func (s *pointsCalculationServiceImpl) ExpireOldPoints(ctx context.Context) error {
	// This would need additional implementation in the repository
	// For now, return nil
	return nil
}