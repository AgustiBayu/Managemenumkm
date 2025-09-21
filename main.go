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
	db.AutoMigrate(&domain.User{}, &domain.Toko{})
	helper.SeedUsers(db)
	validate := validator.New()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokoRepo := repository.NewTokoRepository(db)

	// Services
	userService := service.NewUserService(userRepo, validate)
	tokoService := service.NewTokoService(tokoRepo, validate)

	// Controllers
	userController := controller.NewUserController(userService, tokoService)
	tokoController := controller.NewTokoController(tokoService)

	router := route.NewRouter(userController, tokoController)
	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
