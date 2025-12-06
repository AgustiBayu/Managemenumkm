package helper

import (
	"Managemenumkm/domain"
	"Managemenumkm/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func DBSeed(db *gorm.DB) {
	seedToko(db)
	seedUsers(db)
	seedCate(db)
	seedProducts(db)
	seedMemberTiers(db)
	seedMembers(db)
}

func seedProducts(db *gorm.DB) {
	products := []domain.Product{
		{Name: "Indomie Goreng", SKU: "IDO-GRG-001", Price: 3000, CategoryID: 1, ImageURL: "/static/image/123.png"},
		{Name: "Susu Ultra Coklat", SKU: "ULT-CKL-250", Price: 7000, CategoryID: 1, ImageURL: "/static/image/2.jpeg"},
		{Name: "Teh Pucuk Harum", SKU: "TPH-ORI-350", Price: 3500, CategoryID: 1, ImageURL: "/static/image/3.jpeg"},
		{Name: "Chitato Sapi Panggang", SKU: "CHT-SPG-068", Price: 12000, CategoryID: 1, ImageURL: "/static/image/3.png"},
	}

	// Check if products already exist to avoid duplicates
	var count int64
	db.Model(&domain.Product{}).Count(&count)
	if count == 0 {
		for _, product := range products {
			db.Create(&product)
		}
	}
}

func seedCate(db *gorm.DB) {
	var count int64
	db.Model(&domain.ProductCategory{}).Count(&count)
	if count > 0 {
		return
	}

	cates := []domain.ProductCategory{
		{Category: "Makanan"},
	}
	for _, cate := range cates {
		db.Create(&cate)
	}
}
func seedToko(db *gorm.DB) {
	var count int64
	db.Model(&domain.Toko{}).Count(&count)
	if count > 0 {
		return
	}

	tokos := []domain.Toko{
		{Name: "Mahameru 02", Address: "Sumberayu jln Untung Suropati. blok 02 Pasar Sapii"},
	}
	for _, toko := range tokos {
		db.Create(&toko)
	}
}

func seedUsers(db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	_, err := userRepository.FindByEmail("agustibayusamudro27@gmail.com")
	if err == nil {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := []domain.User{
		{FName: "Agusti", LName: "Bayu Samudro", Email: "agustibayusamudro27@gmail.com", Password: string(hashedPassword), Alamat: "jln. Haji Ahmad Yani RT.RW 004.006 Banyuwangi",
			Thumbnail: "", Number: "081330654123", Role: "superadmin", TokoID: 1},
		{FName: "Rosyita", LName: "Dewi", Email: "noorrosyitadewi027@gmail.com", Password: string(hashedPassword), Alamat: "jln. Haji Ahmad Yani RT.RW 004.006 Kediri",
			Thumbnail: "", Number: "081330100999", Role: "admin", TokoID: 1},
		{FName: "Reka", LName: "Nanda Putri", Email: "rekananputri@gmail.com", Password: string(hashedPassword), Alamat: "jln. Setya Budi RT.RW 005.006 Banyuwangi",
			Thumbnail: "", Number: "081330595025", Role: "cashier", TokoID: 1},
	}

	for _, user := range users {
		db.Create(&user)
	}
}

func seedMemberTiers(db *gorm.DB) {
	var count int64
	db.Model(&domain.MemberTier{}).Count(&count)
	if count > 0 {
		return
	}

	// Create default member tiers for TokoID = 1
	memberTiers := []domain.MemberTier{
		{
			TokoID:       1,
			Name:         "Bronze",
			Description:  "Basic member tier - Welcome to the club!",
			MinPoints:    0,
			MinSpent:     0,
			DiscountRate: 0.00, // 0% discount
			PointRate:    1.0,  // 1 point per Rp 100
			IsDefault:    true,
			IsActive:     true,
		},
		{
			TokoID:       1,
			Name:         "Silver",
			Description:  "Silver member - Enjoy 5% discount!",
			MinPoints:    1000,
			MinSpent:     500000, // Rp 500,000
			DiscountRate: 0.05,   // 5% discount
			PointRate:    1.2,    // 1.2 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
		{
			TokoID:       1,
			Name:         "Gold",
			Description:  "Gold member - Enjoy 10% discount!",
			MinPoints:    5000,
			MinSpent:     2000000, // Rp 2,000,000
			DiscountRate: 0.10,    // 10% discount
			PointRate:    1.5,     // 1.5 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
		{
			TokoID:       1,
			Name:         "Platinum",
			Description:  "Platinum member - Enjoy 15% discount!",
			MinPoints:    10000,
			MinSpent:     5000000, // Rp 5,000,000
			DiscountRate: 0.15,    // 15% discount
			PointRate:    2.0,     // 2 points per Rp 100
			IsDefault:    false,
			IsActive:     true,
		},
	}

	for _, tier := range memberTiers {
		db.Create(&tier)
	}
}

func seedMembers(db *gorm.DB) {
	var count int64
	db.Model(&domain.Member{}).Count(&count)
	if count > 0 {
		return
	}

	// Get the first member tier (Bronze) as default
	var bronzeTier domain.MemberTier
	db.Where("name = ? AND toko_id = ?", "Bronze", uint(1)).First(&bronzeTier)

	// Create sample members
	members := []domain.Member{
		{
			TokoID:        1,
			MemberCode:    "MBR001",
			Name:          "Ahmad Rizki",
			Email:         "ahmad@email.com",
			Phone:         "081234567890",
			Address:       "Jakarta Selatan",
			Birthday:      "1990-01-15",
			Gender:        "Male",
			Notes:         "Pelanggan setia",
			MemberTierID:  bronzeTier.ID,
			TotalPoints:   150,
			TotalSpent:    2500000,
			TotalTransactions: 15,
			Status:        "active",
			IsActive:      true,
			JoinedDate:    time.Now().AddDate(-2, 0, 0), // 2 years ago
		},
		{
			TokoID:        1,
			MemberCode:    "MBR002",
			Name:          "Siti Nurhaliza",
			Email:         "siti@email.com",
			Phone:         "082234567891",
			Address:       "Bandung",
			Birthday:      "1985-05-20",
			Gender:        "Female",
			Notes:         "Pelanggan premium",
			MemberTierID:  bronzeTier.ID,
			TotalPoints:   750,
			TotalSpent:    8000000,
			TotalTransactions: 45,
			Status:        "active",
			IsActive:      true,
			JoinedDate:    time.Now().AddDate(-1, -6, 0), // 1.5 years ago
		},
		{
			TokoID:        1,
			MemberCode:    "MBR003",
			Name:          "Budi Santoso",
			Email:         "budi@email.com",
			Phone:         "08334567892",
			Address:       "Surabaya",
			Birthday:      "1992-08-10",
			Gender:        "Male",
			Notes:         "Pelanggan menengah",
			MemberTierID:  bronzeTier.ID,
			TotalPoints:   320,
			TotalSpent:    4500000,
			TotalTransactions: 28,
			Status:        "active",
			IsActive:      true,
			JoinedDate:    time.Now().AddDate(-1, 0, 0), // 1 year ago
		},
		{
			TokoID:        1,
			MemberCode:    "MBR004",
			Name:          "Dewi Lestari",
			Email:         "dewi@email.com",
			Phone:         "0844567893",
			Address:       "Yogyakarta",
			Birthday:      "1988-12-25",
			Gender:        "Female",
			Notes:         "Pelanggan regular",
			MemberTierID:  bronzeTier.ID,
			TotalPoints:   85,
			TotalSpent:    1200000,
			TotalTransactions: 12,
			Status:        "active",
			IsActive:      true,
			JoinedDate:    time.Now().AddDate(-0, -8, 0), // 8 months ago
		},
		{
			TokoID:        1,
			MemberCode:    "MBR005",
			Name:          "Rudi Hermawan",
			Email:         "rudi@email.com",
			Phone:         "0855678904",
			Address:       "Semarang",
			Birthday:      "1995-03-18",
			Gender:        "Male",
			Notes:         "Pelanggan baru",
			MemberTierID:  bronzeTier.ID,
			TotalPoints:   25,
			TotalSpent:    550000,
			TotalTransactions: 6,
			Status:        "active",
			IsActive:      true,
			JoinedDate:    time.Now().AddDate(-0, -3, 0), // 3 months ago
		},
	}

	for _, member := range members {
		db.Create(&member)
	}
}
