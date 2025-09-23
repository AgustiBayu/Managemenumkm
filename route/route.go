package route

import (
	"Managemenumkm/controller"
	"Managemenumkm/domain"
	"Managemenumkm/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(userController controller.UserController, tokoController controller.TokoController) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// --- Public Routes ---
	router.GET("/login", userController.ShowLoginForm)
	router.POST("/login", userController.Login)
	router.GET("/logout", userController.Logout)

	// --- Authenticated Routes (All Roles) ---
	router.GET("/dashboard", middleware.Authorize(userController.Dashboard, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.GET("/profile", middleware.Authorize(userController.ShowProfile, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))

	// --- User Management (Admin & Superadmin) ---
	router.GET("/user", middleware.Authorize(userController.FindAll, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.GET("/user/add", middleware.Authorize(userController.Create, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.POST("/user/add", middleware.Authorize(userController.Create, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.GET("/user/edit/:userID", middleware.Authorize(userController.FindById, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.PUT("/user/edit/:userID", middleware.Authorize(userController.Update, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.DELETE("/user/delete/:userID", middleware.Authorize(userController.Delete, domain.RoleSuperAdmin, domain.RoleAdmin))

	// --- Toko Management ---
	// Admins can see the list of all stores, but can't act on them (this can be refined in the view)
	router.GET("/toko", middleware.Authorize(tokoController.FindAll, domain.RoleSuperAdmin, domain.RoleAdmin))

	// Only Superadmin can Create, Update, Delete stores
	router.GET("/toko/add", middleware.Authorize(tokoController.Create, domain.RoleSuperAdmin))
	router.POST("/toko/add", middleware.Authorize(tokoController.Create, domain.RoleSuperAdmin))
	router.GET("/toko/edit/:tokoID", middleware.Authorize(tokoController.FindById, domain.RoleSuperAdmin))
	router.POST("/toko/edit/:tokoID", middleware.Authorize(tokoController.Update, domain.RoleSuperAdmin))
	router.GET("/toko/delete/:tokoID", middleware.Authorize(tokoController.Delete, domain.RoleSuperAdmin))

	return router
}
