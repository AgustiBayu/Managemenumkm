package service

import (
	"Managemenumkm/domain"
	"Managemenumkm/repository"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type MemberServiceImpl struct {
	MemberRepository     repository.MemberRepository
	MemberTierRepository repository.MemberTierRepository
	Validate             *validator.Validate
}

func NewMemberService(memberRepository repository.MemberRepository, memberTierRepository repository.MemberTierRepository, validate *validator.Validate) MemberService {
	return &MemberServiceImpl{
		MemberRepository:     memberRepository,
		MemberTierRepository: memberTierRepository,
		Validate:             validate,
	}
}

func (s *MemberServiceImpl) CreateMember(request domain.Member) (domain.Member, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return domain.Member{}, err
	}

	// Generate unique member code
	memberCode := s.generateMemberCode()

	// Get default tier for new members
	defaultTier, err := s.MemberTierRepository.GetDefaultTier(request.TokoID)
	if err != nil {
		return domain.Member{}, err
	}

	member := domain.Member{
		TokoID:       request.TokoID,
		MemberCode:   memberCode,
		Name:         request.Name,
		Email:        request.Email,
		Phone:        request.Phone,
		Address:      request.Address,
		MemberTierID: defaultTier.ID,
		TotalPoints:  0,
		TotalSpent:   0,
		IsActive:     true,
		JoinedDate:   time.Now(),
	}

	newMember, err := s.MemberRepository.Create(member)
	if err != nil {
		return domain.Member{}, err
	}

	return newMember, nil
}

func (s *MemberServiceImpl) UpdateMember(memberID uint, request domain.Member) (domain.Member, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		return domain.Member{}, err
	}

	existingMember, err := s.MemberRepository.FindByID(memberID)
	if err != nil {
		return domain.Member{}, err
	}

	// Debug logging
	log.Printf("Updating member %d. Current tier ID: %d, Request tier ID: %d",
		memberID, existingMember.MemberTierID, request.MemberTierID)

	existingMember.Name = request.Name
	existingMember.Email = request.Email
	existingMember.Phone = request.Phone
	existingMember.Address = request.Address
	existingMember.Birthday = request.Birthday
	existingMember.Gender = request.Gender
	existingMember.Notes = request.Notes
	existingMember.Status = request.Status

	// Only update MemberTierID if it's provided (not 0)
	if request.MemberTierID != 0 {
		log.Printf("Updating MemberTierID from %d to %d", existingMember.MemberTierID, request.MemberTierID)
		existingMember.MemberTierID = request.MemberTierID
	} else {
		log.Printf("MemberTierID is 0 in request, keeping existing value: %d", existingMember.MemberTierID)
	}

	updatedMember, err := s.MemberRepository.Update(existingMember)
	if err != nil {
		return domain.Member{}, err
	}

	log.Printf("Successfully updated member. Final tier ID: %d", updatedMember.MemberTierID)
	return updatedMember, nil
}

func (s *MemberServiceImpl) DeleteMember(memberID uint) error {
	err := s.MemberRepository.Delete(memberID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MemberServiceImpl) GetMemberByID(memberID uint) (domain.Member, error) {
	member, err := s.MemberRepository.FindByID(memberID)
	if err != nil {
		return domain.Member{}, err
	}
	return member, nil
}

func (s *MemberServiceImpl) GetMembersByTokoID(tokoID uint) ([]domain.Member, error) {
	members, err := s.MemberRepository.FindByTokoID(tokoID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *MemberServiceImpl) GetMembersWithFilter(filter domain.MemberFilter) ([]domain.Member, error) {
	members, err := s.MemberRepository.FindWithFilter(filter)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *MemberServiceImpl) FindByMemberCode(memberCode string, tokoID uint) (domain.Member, error) {
	member, err := s.MemberRepository.FindByMemberCode(memberCode, tokoID)
	if err != nil {
		return domain.Member{}, err
	}
	return member, nil
}

func (s *MemberServiceImpl) FindByPhone(phone string, tokoID uint) (domain.Member, error) {
	member, err := s.MemberRepository.FindByPhone(phone, tokoID)
	if err != nil {
		return domain.Member{}, err
	}
	return member, nil
}

func (s *MemberServiceImpl) EarnPoints(memberID uint, transactionID uint, amount float64, description string) (domain.MemberTransaction, error) {
	member, err := s.MemberRepository.FindByID(memberID)
	if err != nil {
		return domain.MemberTransaction{}, err
	}

	if !member.IsActive {
		return domain.MemberTransaction{}, errors.New("member is not active")
	}

	// Calculate points (1 point per Rp 100 spent)
	pointsEarned := int(amount) / 100

	if pointsEarned <= 0 {
		return domain.MemberTransaction{}, errors.New("no points earned for this transaction")
	}

	balanceBefore := member.TotalPoints
	balanceAfter := balanceBefore + pointsEarned

	// Create member transaction
	memberTransaction := domain.MemberTransaction{
		TokoID:        member.TokoID,
		MemberID:      memberID,
		TransactionID: &transactionID,
		Type:          "EARN",
		Points:        pointsEarned,
		Amount:        amount,
		Description:   description,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}

	// Update member points and total spent
	member.TotalPoints = balanceAfter
	member.TotalSpent += amount
	now := time.Now()
	member.LastVisitDate = &now

	// Save member transaction (you'll need to create this repository method)
	// For now, let's skip the actual saving to DB and just update member points
	err = s.MemberRepository.UpdatePoints(memberID, balanceAfter)
	if err != nil {
		return domain.MemberTransaction{}, err
	}

	// Update member tier if needed
	_, err = s.UpdateMemberTier(memberID)
	if err != nil {
		log.Printf("Error updating member tier: %v", err)
		// Continue even if tier update fails
	}

	return memberTransaction, nil
}

func (s *MemberServiceImpl) RedeemPoints(memberID uint, points int, description string) (domain.MemberTransaction, error) {
	member, err := s.MemberRepository.FindByID(memberID)
	if err != nil {
		return domain.MemberTransaction{}, err
	}

	if !member.IsActive {
		return domain.MemberTransaction{}, errors.New("member is not active")
	}

	if member.TotalPoints < points {
		return domain.MemberTransaction{}, errors.New("insufficient points")
	}

	balanceBefore := member.TotalPoints
	balanceAfter := balanceBefore - points

	// Create member transaction
	memberTransaction := domain.MemberTransaction{
		TokoID:        member.TokoID,
		MemberID:      memberID,
		Type:          "REDEEM",
		Points:        -points,
		Description:   description,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}

	// Update member points
	err = s.MemberRepository.UpdatePoints(memberID, balanceAfter)
	if err != nil {
		return domain.MemberTransaction{}, err
	}

	return memberTransaction, nil
}

func (s *MemberServiceImpl) GetMemberTransactions(memberID uint, limit int, offset int) ([]domain.MemberTransaction, error) {
	transactions, err := s.MemberRepository.GetMemberTransactions(memberID, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *MemberServiceImpl) GetPointsBalance(memberID uint) (int, error) {
	balance, err := s.MemberRepository.GetPointsBalance(memberID)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *MemberServiceImpl) UpdateMemberTier(memberID uint) (domain.Member, error) {
	member, err := s.MemberRepository.FindByID(memberID)
	if err != nil {
		return domain.Member{}, err
	}

	// Find appropriate tier based on points and spent amount
	// First check by points
	tierByPoints, err := s.MemberTierRepository.FindByPoints(member.TotalPoints, member.TokoID)
	if err != nil {
		tierByPoints, err = s.MemberTierRepository.GetDefaultTier(member.TokoID)
		if err != nil {
			return domain.Member{}, err
		}
	}

	// Also check by spent amount and pick the better tier
	tierBySpent, err := s.MemberTierRepository.FindBySpentAmount(member.TotalSpent, member.TokoID)
	if err != nil {
		tierBySpent = tierByPoints
	}

	// Choose the tier with better benefits (higher discount rate or point rate)
	selectedTier := tierByPoints
	if tierBySpent.DiscountRate > tierByPoints.DiscountRate || tierBySpent.PointRate > tierByPoints.PointRate {
		selectedTier = tierBySpent
	}

	// Update member tier if it's different
	if member.MemberTierID != selectedTier.ID {
		member.MemberTierID = selectedTier.ID
		updatedMember, err := s.MemberRepository.Update(member)
		if err != nil {
			return domain.Member{}, err
		}
		return updatedMember, nil
	}

	return member, nil
}

func (s *MemberServiceImpl) generateMemberCode() string {
	// Generate a unique member code in format: MBR + 8 digit random number
	// Example: MBR12345678
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(99999999) // 0 to 99,999,999
	return fmt.Sprintf("MBR%08d", randomNum)
}

func (s *MemberServiceImpl) GetMemberTiers(tokoID uint) ([]domain.MemberTier, error) {
	return s.MemberTierRepository.FindByTokoID(tokoID)
}

// AutoInactiveMembers marks members as inactive if they haven't made purchases in 6 months
func (s *MemberServiceImpl) AutoInactiveMembers(tokoID uint) error {
	// Get all active members
	members, err := s.MemberRepository.FindActiveMembers(tokoID)
	if err != nil {
		return err
	}

	// Calculate 6 months ago from today
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	// Check each member's last activity
	for _, member := range members {
		shouldInactive := false

		// Check last visit date
		if member.LastVisitDate != nil {
			if member.LastVisitDate.Before(sixMonthsAgo) {
				shouldInactive = true
			}
		} else {
			// If no last visit date, check joined date
			if member.JoinedDate.Before(sixMonthsAgo) {
				shouldInactive = true
			}
		}

		// If member should be inactive, update status
		if shouldInactive && member.Status == "active" {
			member.Status = "inactive"
			member.IsActive = false
			_, err := s.MemberRepository.Update(member)
			if err != nil {
				log.Printf("Failed to update member %s to inactive: %v", member.MemberCode, err)
			} else {
				log.Printf("Member %s marked as inactive due to 6 months inactivity", member.MemberCode)
			}
		}
	}

	return nil
}

// Helper function to clean phone number
func cleanPhoneNumber(phone string) string {
	// Remove all non-digit characters
	phone = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Remove leading 62 or 0 and replace with 62 if needed
	if strings.HasPrefix(phone, "62") {
		return phone
	}
	if strings.HasPrefix(phone, "0") {
		return "62" + phone[1:]
	}
	return phone
}

func (s *MemberServiceImpl) ImportMembers(members []ImportMember, tokoID uint, skipDuplicates, updateExisting bool) (ImportMembersResponse, error) {
	response := ImportMembersResponse{}

	// Get member tiers for tier mapping
	memberTiers, err := s.MemberTierRepository.FindByTokoID(tokoID)
	if err != nil {
		return response, fmt.Errorf("failed to get member tiers: %w", err)
	}

	// Create a map for tier name to ID
	tierMap := make(map[string]uint)
	for _, tier := range memberTiers {
		tierMap[strings.ToLower(tier.Name)] = tier.ID
	}

	// Get default tier
	defaultTier, err := s.MemberTierRepository.GetDefaultTier(tokoID)
	if err != nil {
		return response, fmt.Errorf("failed to get default member tier: %w", err)
	}

	for _, importMember := range members {
		// Check if member already exists by phone
		existingMember, err := s.MemberRepository.FindByPhone(importMember.Phone, tokoID)
		if err == nil {
			// Member exists
			if skipDuplicates && !updateExisting {
				response.Skipped++
				continue
			}

			if updateExisting {
				// Update existing member
				existingMember.Name = importMember.Name
				existingMember.Email = importMember.Email
				existingMember.Address = importMember.Address
				existingMember.Birthday = importMember.Birthday
				existingMember.Gender = importMember.Gender
				existingMember.Notes = importMember.Notes

				// Update status if provided
				if importMember.Status != "" {
					existingMember.Status = importMember.Status
				}

				// Update tier if provided and valid
				if importMember.Tier != "" {
					if tierID, exists := tierMap[strings.ToLower(importMember.Tier)]; exists {
						existingMember.MemberTierID = tierID
					}
				}

				_, err = s.MemberRepository.Update(existingMember)
				if err != nil {
					response.Errors++
					log.Printf("Failed to update existing member with phone %s: %v", importMember.Phone, err)
					continue
				}
				response.Imported++
			} else {
				response.Skipped++
			}
		} else {
			// Member doesn't exist, create new one
			memberCode := s.generateMemberCode()

			// Determine tier ID
			var memberTierID uint = defaultTier.ID
			if importMember.Tier != "" {
				if tierID, exists := tierMap[strings.ToLower(importMember.Tier)]; exists {
					memberTierID = tierID
				}
			}

			// Determine status
			status := "active"
			if importMember.Status != "" {
				status = importMember.Status
			}

			member := domain.Member{
				TokoID:       tokoID,
				MemberCode:   memberCode,
				Name:         importMember.Name,
				Email:        importMember.Email,
				Phone:        importMember.Phone,
				Address:      importMember.Address,
				Birthday:     importMember.Birthday,
				Gender:       importMember.Gender,
				MemberTierID: memberTierID,
				TotalPoints:  0,
				TotalSpent:   0,
				IsActive:     true,
				Status:       status,
				Notes:        importMember.Notes,
				JoinedDate:   time.Now(),
			}

			_, err = s.MemberRepository.Create(member)
			if err != nil {
				response.Errors++
				log.Printf("Failed to create member with phone %s: %v", importMember.Phone, err)
				continue
			}
			response.Imported++
		}
	}

	return response, nil
}
