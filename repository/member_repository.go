package repository

import (
	"Managemenumkm/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member domain.Member) (domain.Member, error)
	Update(member domain.Member) (domain.Member, error)
	Delete(memberID uint) error
	FindByID(memberID uint) (domain.Member, error)
	FindByTokoID(tokoID uint) ([]domain.Member, error)
	FindActiveMembers(tokoID uint) ([]domain.Member, error)
	FindByMemberCode(memberCode string, tokoID uint) (domain.Member, error)
	FindByPhone(phone string, tokoID uint) (domain.Member, error)
	FindWithFilter(filter domain.MemberFilter) ([]domain.Member, error)
	UpdatePoints(memberID uint, points int) error
	GetMemberTransactions(memberID uint, limit int, offset int) ([]domain.MemberTransaction, error)
	GetPointsBalance(memberID uint) (int, error)
}

type MemberRepositoryImpl struct {
	DB *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &MemberRepositoryImpl{DB: db}
}

func (r *MemberRepositoryImpl) Create(member domain.Member) (domain.Member, error) {
	err := r.DB.Create(&member).Error
	return member, err
}

func (r *MemberRepositoryImpl) Update(member domain.Member) (domain.Member, error) {
	// Debug logging to help troubleshoot update issues
	fmt.Printf("Repository: Updating member %d with MemberTierID: %d\n", member.ID, member.MemberTierID)

	// Try using a more explicit update approach
	result := r.DB.Model(&domain.Member{}).Where("id = ?", member.ID).Updates(map[string]interface{}{
		"name":           member.Name,
		"email":          member.Email,
		"phone":          member.Phone,
		"address":        member.Address,
		"birthday":       member.Birthday,
		"gender":         member.Gender,
		"notes":          member.Notes,
		"status":         member.Status,
		"member_tier_id": member.MemberTierID,
		"updated_at":     time.Now(),
	})

	if result.Error != nil {
		fmt.Printf("Repository: Error updating member with Updates: %v\n", result.Error)
		return member, result.Error
	}

	if result.RowsAffected == 0 {
		fmt.Printf("Repository: No rows affected during update\n")
	} else {
		fmt.Printf("Repository: Successfully updated %d rows\n", result.RowsAffected)
	}

	// Verify the update by fetching the member again
	updatedMember, err := r.FindByID(member.ID)
	if err != nil {
		fmt.Printf("Repository: Error fetching updated member: %v\n", err)
		return member, err
	}

	fmt.Printf("Repository: Verified updated member %d has MemberTierID: %d\n", updatedMember.ID, updatedMember.MemberTierID)
	return updatedMember, nil
}

func (r *MemberRepositoryImpl) Delete(memberID uint) error {
	err := r.DB.Delete(&domain.Member{}, memberID).Error
	return err
}

func (r *MemberRepositoryImpl) FindByID(memberID uint) (domain.Member, error) {
	var member domain.Member
	err := r.DB.Preload("MemberTier").First(&member, memberID).Error
	return member, err
}

func (r *MemberRepositoryImpl) FindByTokoID(tokoID uint) ([]domain.Member, error) {
	var members []domain.Member
	err := r.DB.Preload("MemberTier").Where("toko_id = ?", tokoID).Find(&members).Error
	return members, err
}

func (r *MemberRepositoryImpl) FindByMemberCode(memberCode string, tokoID uint) (domain.Member, error) {
	var member domain.Member
	err := r.DB.Preload("MemberTier").Where("member_code = ? AND toko_id = ?", memberCode, tokoID).First(&member).Error
	return member, err
}

func (r *MemberRepositoryImpl) FindByPhone(phone string, tokoID uint) (domain.Member, error) {
	var member domain.Member
	err := r.DB.Preload("MemberTier").Where("phone = ? AND toko_id = ?", phone, tokoID).First(&member).Error
	return member, err
}

func (r *MemberRepositoryImpl) UpdatePoints(memberID uint, points int) error {
	err := r.DB.Model(&domain.Member{}).Where("id = ?", memberID).Update("total_points", points).Error
	return err
}

func (r *MemberRepositoryImpl) GetMemberTransactions(memberID uint, limit int, offset int) ([]domain.MemberTransaction, error) {
	var transactions []domain.MemberTransaction
	err := r.DB.Where("member_id = ?", memberID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *MemberRepositoryImpl) GetPointsBalance(memberID uint) (int, error) {
	var member domain.Member
	err := r.DB.Select("total_points").First(&member, memberID).Error
	return member.TotalPoints, err
}

func (r *MemberRepositoryImpl) FindActiveMembers(tokoID uint) ([]domain.Member, error) {
	var members []domain.Member
	err := r.DB.Preload("MemberTier").Where("toko_id = ? AND status = ?", tokoID, "active").Find(&members).Error
	return members, err
}

func (r *MemberRepositoryImpl) FindWithFilter(filter domain.MemberFilter) ([]domain.Member, error) {
	var members []domain.Member
	query := r.DB.Preload("MemberTier").Table("members").Where("members.toko_id = ?", filter.TokoID)

	// Add search filter if provided
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("(members.name LIKE ? OR members.phone LIKE ? OR members.member_code LIKE ? OR members.email LIKE ?)",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Add tier filter if provided
	if filter.Tier != "" {
		query = query.Joins("JOIN member_tiers ON members.member_tier_id = member_tiers.id").
			Where("LOWER(member_tiers.name) = LOWER(?)", filter.Tier)
	}

	// Add status filter if provided
	if filter.Status != "" {
		query = query.Where("LOWER(members.status) = LOWER(?)", filter.Status)
	}

	err := query.Find(&members).Error
	return members, err
}
