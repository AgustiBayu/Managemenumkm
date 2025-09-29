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
	db.AutoMigrate(&domain.User{}, &domain.Toko{}, &domain.ProductCategory{}, &domain.Product{})
	helper.DBSeed(db)
	validate := validator.New()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	cateRepo := repository.NewProductCategoryRepository(db)
	proRepo := repository.NewProductRepository(db)

	// Services
	userService := service.NewUserService(userRepo, validate)
	tokoService := service.NewTokoService(tokoRepo, validate)
	cateService := service.NewProductCategoryService(cateRepo, validate)
	proService := service.NewProductService(proRepo, cateRepo, validate)

	// Controllers
	userController := controller.NewUserController(userService, tokoService)
	tokoController := controller.NewTokoController(tokoService)
	cateController := controller.NewProductCategoryController(cateService)
	proController := controller.NewProductController(proService, cateService)

	router := route.NewRouter(userController, tokoController, cateController, proController)
	println("Server running at https://localhost:8443")
	http.ListenAndServeTLS(":8443", "certs/server.crt", "certs/server.key", router)
}
