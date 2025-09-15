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
)

func main() {

	helper.InitJWT()
	db := app.DB()
	db.AutoMigrate(&domain.User{})
	helper.SeedUsers(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo)

	// Controllers
	userController := controller.NewUserController(userService)

	router := route.NewRouter(userController)
	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
