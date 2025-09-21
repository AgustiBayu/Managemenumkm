package helper

import (
	"Managemenumkm/domain"
	"Managemenumkm/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	_, err := userRepository.FindByEmail("agustibayusamudro27@gmail.com")
	if err == nil {
		// User already exists
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := []domain.User{
		{FName: "Agusti", LName: "Bayu Samudro", Email: "agustibayusamudro27@gmail.com", Password: string(hashedPassword), Alamat: "jln. Haji Ahmad Yani RT.RW 004.006 Banyuwangi",
			Thumbnail: "", Number: "081330654123", Role: "admin"},
		{FName: "Reka", LName: "Nanda Putri", Email: "rekananputri@gmail.com", Password: string(hashedPassword), Alamat: "jln. Setya Budi RT.RW 005.006 Banyuwangi",
			Thumbnail: "", Number: "081330595025", Role: "cashier"},
	}

	for _, user := range users {
		db.Create(&user)
	}
}
