package repository

import (
	"Managemenumkm/domain"
	"time"
	"gorm.io/gorm"
)

type DiscountRepositoryImpl struct {
	DB *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &DiscountRepositoryImpl{DB: db}
}

// Basic CRUD operations
func (r *DiscountRepositoryImpl) Create(discount domain.Discount) (domain.Discount, error) {
	err := r.DB.Create(&discount).Error
	return discount, err
}

func (r *DiscountRepositoryImpl) Update(discount domain.Discount) (domain.Discount, error) {
	err := r.DB.Save(&discount).Error
	return discount, err
}

func (r *DiscountRepositoryImpl) Delete(discountID uint) error {
	err := r.DB.Delete(&domain.Discount{}, discountID).Error
	return err
}

func (r *DiscountRepositoryImpl) FindByID(discountID uint) (domain.Discount, error) {
	var discount domain.Discount
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		First(&discount, discountID).Error
	return discount, err
}

func (r *DiscountRepositoryImpl) FindByTokoID(tokoID uint) ([]domain.Discount, error) {
	var discounts []domain.Discount
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ?", tokoID).
		Order("created_at DESC").
		Find(&discounts).Error
	return discounts, err
}

// Advanced queries
func (r *DiscountRepositoryImpl) FindActiveDiscounts(tokoID uint) ([]domain.Discount, error) {
	var discounts []domain.Discount
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ? AND is_active = ?", tokoID, true).
		Order("created_at DESC").
		Find(&discounts).Error
	return discounts, err
}

func (r *DiscountRepositoryImpl) FindValidDiscounts(tokoID uint, currentDate time.Time) ([]domain.Discount, error) {
	var discounts []domain.Discount
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND (usage_limit = 0 OR usage_count < usage_limit)",
			tokoID, true, currentDate, currentDate).
		Order("created_at DESC").
		Find(&discounts).Error
	return discounts, err
}

func (r *DiscountRepositoryImpl) FindWithFilter(filter domain.DiscountFilter) ([]domain.Discount, error) {
	var discounts []domain.Discount
	query := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ?", filter.TokoID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.ApplicableTo != "" {
		query = query.Where("applicable_to = ?", filter.ApplicableTo)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("(name LIKE ? OR description LIKE ?)", searchTerm, searchTerm)
	}

	err := query.Order("created_at DESC").Find(&discounts).Error
	return discounts, err
}

// Discount usage tracking
func (r *DiscountRepositoryImpl) GetDiscountUsage(discountID uint) (int, error) {
	var discount domain.Discount
	err := r.DB.Select("usage_count").First(&discount, discountID).Error
	return discount.UsageCount, err
}

func (r *DiscountRepositoryImpl) IncrementUsage(discountID uint) error {
	err := r.DB.Model(&domain.Discount{}).
		Where("id = ?", discountID).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
	return err
}

// Special offers
func (r *DiscountRepositoryImpl) CreateSpecialOffer(specialOffer domain.SpecialOffer) (domain.SpecialOffer, error) {
	err := r.DB.Create(&specialOffer).Error
	return specialOffer, err
}

func (r *DiscountRepositoryImpl) UpdateSpecialOffer(specialOffer domain.SpecialOffer) (domain.SpecialOffer, error) {
	err := r.DB.Save(&specialOffer).Error
	return specialOffer, err
}

func (r *DiscountRepositoryImpl) DeleteSpecialOffer(specialOfferID uint) error {
	err := r.DB.Delete(&domain.SpecialOffer{}, specialOfferID).Error
	return err
}

func (r *DiscountRepositoryImpl) FindActiveSpecialOffers(tokoID uint) ([]domain.SpecialOffer, error) {
	var specialOffers []domain.SpecialOffer
	currentTime := time.Now()
	err := r.DB.Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ?",
		tokoID, true, currentTime, currentTime).
		Order("priority DESC, created_at DESC").
		Find(&specialOffers).Error
	return specialOffers, err
}

func (r *DiscountRepositoryImpl) FindSpecialOfferProducts(specialOfferID uint) ([]domain.SpecialOfferProduct, error) {
	var products []domain.SpecialOfferProduct
	err := r.DB.Preload("Product").
		Where("special_offer_id = ?", specialOfferID).
		Find(&products).Error
	return products, err
}

// Flash sales
func (r *DiscountRepositoryImpl) CreateFlashSale(flashSale domain.FlashSale) (domain.FlashSale, error) {
	err := r.DB.Create(&flashSale).Error
	return flashSale, err
}

func (r *DiscountRepositoryImpl) UpdateFlashSale(flashSale domain.FlashSale) (domain.FlashSale, error) {
	err := r.DB.Save(&flashSale).Error
	return flashSale, err
}

func (r *DiscountRepositoryImpl) DeleteFlashSale(flashSaleID uint) error {
	err := r.DB.Delete(&domain.FlashSale{}, flashSaleID).Error
	return err
}

func (r *DiscountRepositoryImpl) FindActiveFlashSales(tokoID uint) ([]domain.FlashSale, error) {
	var flashSales []domain.FlashSale
	currentTime := time.Now()
	err := r.DB.Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ?",
		tokoID, true, currentTime, currentTime).
		Order("created_at DESC").
		Find(&flashSales).Error
	return flashSales, err
}

func (r *DiscountRepositoryImpl) FindFlashSaleProducts(flashSaleID uint) ([]domain.FlashSaleProduct, error) {
	var products []domain.FlashSaleProduct
	err := r.DB.Preload("Product").
		Where("flash_sale_id = ?", flashSaleID).
		Find(&products).Error
	return products, err
}

// Discount usage tracking
func (r *DiscountRepositoryImpl) CreateDiscountUsage(usage domain.DiscountUsage) (domain.DiscountUsage, error) {
	err := r.DB.Create(&usage).Error
	return usage, err
}

func (r *DiscountRepositoryImpl) FindDiscountUsagesByTransaction(transactionID uint) ([]domain.DiscountUsage, error) {
	var usages []domain.DiscountUsage
	err := r.DB.Preload("Discount").
		Where("transaction_id = ?", transactionID).
		Find(&usages).Error
	return usages, err
}

func (r *DiscountRepositoryImpl) FindProductDiscounts(productID uint, tokoID uint) ([]domain.Discount, error) {
	var discounts []domain.Discount
	currentTime := time.Now()
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND (usage_limit = 0 OR usage_count < usage_limit) AND (applicable_to = 'ALL' OR applicable_to = 'PRODUCT' AND product_id = ?)",
			tokoID, true, currentTime, currentTime, productID).
		Find(&discounts).Error
	return discounts, err
}

func (r *DiscountRepositoryImpl) FindCategoryDiscounts(categoryID uint, tokoID uint) ([]domain.Discount, error) {
	var discounts []domain.Discount
	currentTime := time.Now()
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND (usage_limit = 0 OR usage_count < usage_limit) AND (applicable_to = 'ALL' OR applicable_to = 'CATEGORY' AND category_id = ?)",
			tokoID, true, currentTime, currentTime, categoryID).
		Find(&discounts).Error
	return discounts, err
}

func (r *DiscountRepositoryImpl) FindMemberTierDiscounts(memberTierID uint, tokoID uint) ([]domain.Discount, error) {
	var discounts []domain.Discount
	currentTime := time.Now()
	err := r.DB.Preload("Product").
		Preload("Category").
		// MemberTier preload removed - tierless membership system
		Where("toko_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND (usage_limit = 0 OR usage_count < usage_limit) AND (applicable_to = 'ALL' OR applicable_to = 'MEMBER_TIER' AND member_tier_id = ?)",
			tokoID, true, currentTime, currentTime, memberTierID).
		Find(&discounts).Error
	return discounts, err
}