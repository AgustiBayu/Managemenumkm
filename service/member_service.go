package service

import (
	"Managemenumkm/domain"
)

type MemberService interface {
	CreateMember(request domain.Member) (domain.Member, error)
	UpdateMember(memberID uint, request domain.Member) (domain.Member, error)
	DeleteMember(memberID uint) error
	GetMemberByID(memberID uint) (domain.Member, error)
	GetMembersByTokoID(tokoID uint) ([]domain.Member, error)
	GetMembersWithFilter(filter domain.MemberFilter) ([]domain.Member, error)
	FindByMemberCode(memberCode string, tokoID uint) (domain.Member, error)
	FindByPhone(phone string, tokoID uint) (domain.Member, error)
	EarnPoints(memberID uint, transactionID uint, amount float64, description string) (domain.MemberTransaction, error)
	RedeemPoints(memberID uint, points int, description string) (domain.MemberTransaction, error)
	GetMemberTransactions(memberID uint, limit int, offset int) ([]domain.MemberTransaction, error)
	GetPointsBalance(memberID uint) (int, error)
	// UpdateMemberTier and GetMemberTiers removed - tierless membership system
	AutoInactiveMembers(tokoID uint) error
	ImportMembers(members []ImportMember, tokoID uint, skipDuplicates, updateExisting bool) (ImportMembersResponse, error)
}

type ImportMember struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Birthday string `json:"birthday"`
	Gender  string `json:"gender"`
	Status  string `json:"status"`
	Notes   string `json:"notes"`
}

type ImportMembersResponse struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
	Message  string `json:"message,omitempty"`
}

type MemberServiceRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"email"`
	Phone   string `json:"phone" validate:"required"`
	Address string `json:"address"`
}
