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
}
func seedCate(db *gorm.DB) {
	cates := []domain.ProductCategory{
		{Category: "Makanan"},
	}
	for _, cate := range cates {
		db.Create(&cate)
	}
}
func seedToko(db *gorm.DB) {
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
