package main

import (
	"Managemenumkm/app"
	"Managemenumkm/controller"
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/repository"
	"Managemenumkm/route"
	"Managemenumkm/service"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func main() {

	helper.InitJWT()
	db := app.DB()
	db.AutoMigrate(&domain.User{}, &domain.Toko{}, &domain.ProductCategory{}, &domain.Product{}, &domain.ProductBatch{}, &domain.Transaction{}, &domain.TransactionItem{}, &domain.MemberTier{}, &domain.Member{}, &domain.MemberTransaction{}, &domain.PointRedemptionRule{})
	helper.DBSeed(db)
	validate := validator.New()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	cateRepo := repository.NewProductCategoryRepository(db)
	proRepo := repository.NewProductRepository(db)
	batchRepo := repository.NewProductBatchRepository(db)
	transRepo := repository.NewTransactionRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	memberTierRepo := repository.NewMemberTierRepository(db)

	// Services
	userService := service.NewUserService(userRepo, validate)
	tokoService := service.NewTokoService(tokoRepo, validate)
	cateService := service.NewProductCategoryService(cateRepo, validate)
	batchService := service.NewProductBatchService(batchRepo, validate)
	proService := service.NewProductService(proRepo, batchRepo, cateRepo, validate)
	memberService := service.NewMemberService(memberRepo, memberTierRepo, validate)
	transService := service.NewTransactionService(db, transRepo, proRepo, memberService)

	// Controllers
	userController := controller.NewUserController(userService, tokoService)
	tokoController := controller.NewTokoController(tokoService)
	cateController := controller.NewProductCategoryController(cateService)
	proController := controller.NewProductController(proService, batchService, cateService)
	posController := controller.NewPosController(transService, proService)
	memberController := controller.NewMemberController(memberService)

	router := route.NewRouter(userController, tokoController, cateController, proController, posController, memberController)
	println("Server running at https://localhost:8443")

	// For production with HTTPS, use cert and key files
	// For development, you can generate self-signed certificates:
	// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", router)
	if err != nil {
		println("HTTPS server failed to start, falling back to HTTP on :8080")
		println("Error:", err.Error())
		println("Server running at http://localhost:8080")
		http.ListenAndServe(":8080", router)
	}
}
