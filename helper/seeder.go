package helper

import (
	"Managemenumkm/domain"
	"Managemenumkm/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func DBSeed(db *gorm.DB) {
	seedToko(db)
	seedUsers(db)
	seedCate(db)
	seedProducts(db)
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
