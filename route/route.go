package route

import (
	"Managemenumkm/controller"
	"Managemenumkm/domain"
	"Managemenumkm/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(userController controller.UserController, tokoController controller.TokoController, categoryController controller.ProductCategoryController, productController controller.ProductController, posController controller.PosController) *httprouter.Router {
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
	router.GET("/toko", middleware.Authorize(tokoController.FindAll, domain.RoleSuperAdmin, domain.RoleAdmin))
	router.GET("/toko/add", middleware.Authorize(tokoController.Create, domain.RoleSuperAdmin))
	router.POST("/toko/add", middleware.Authorize(tokoController.Create, domain.RoleSuperAdmin))
	router.GET("/toko/edit/:tokoID", middleware.Authorize(tokoController.FindById, domain.RoleSuperAdmin))
	router.POST("/toko/edit/:tokoID", middleware.Authorize(tokoController.Update, domain.RoleSuperAdmin))
	router.GET("/toko/delete/:tokoID", middleware.Authorize(tokoController.Delete, domain.RoleSuperAdmin))

	// --- Category Product Management ---
	router.GET("/category", middleware.Authorize(categoryController.FindAll, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.GET("/category/add", middleware.Authorize(categoryController.Create, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.POST("/category/add", middleware.Authorize(categoryController.Create, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.GET("/category/edit/:categoryID", middleware.Authorize(categoryController.FindById, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.POST("/category/edit/:categoryID", middleware.Authorize(categoryController.Update, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))
	router.GET("/category/delete/:categoryID", middleware.Authorize(categoryController.Delete, domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier))

	// --- Product Management ---
	// Allow all authenticated users to see the product list
	authAll := []domain.Role{domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleCashier}

	// Master Product Routes
	router.GET("/product", middleware.Authorize(productController.FindAll, authAll...))
	router.GET("/product/add", middleware.Authorize(productController.CreateView, authAll...))
	router.POST("/product/add", middleware.Authorize(productController.Create, authAll...))
	router.GET("/product/edit/:productId", middleware.Authorize(productController.FindById, authAll...))
	router.POST("/product/edit/:productId", middleware.Authorize(productController.Update, authAll...))
	router.GET("/product/delete/:productId", middleware.Authorize(productController.Delete, authAll...))

	// Stock and Batch Management Routes
	router.GET("/product/stock/add", middleware.Authorize(productController.AddStockView, authAll...))
	router.POST("/product/stock/scan", middleware.Authorize(productController.AddStockScanBarcode, authAll...))
	router.POST("/product/stock/save", middleware.Authorize(productController.AddStockSaveBatch, authAll...))
	router.GET("/product/batches/:productId", middleware.Authorize(productController.ListBatches, authAll...))
	router.GET("/product/batch/edit/:batchId", middleware.Authorize(productController.EditBatchView, authAll...))
	router.POST("/product/batch/edit/:batchId", middleware.Authorize(productController.EditBatch, authAll...))

	// API route for image upload
	router.POST("/api/products/:productId/image", middleware.Authorize(productController.UploadImage, authAll...))

	// --- POS (Point of Sale) ---
	router.GET("/pos", middleware.Authorize(posController.ShowPosPage, authAll...))
	router.GET("/api/pos/products", middleware.Authorize(posController.GetProducts, authAll...))
	router.POST("/api/pos/transactions", middleware.Authorize(posController.CreateTransaction, authAll...))

	return router
}