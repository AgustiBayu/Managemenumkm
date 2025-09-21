package route

import (
	"Managemenumkm/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Middleware wrapper
func protected(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// This is a simplified middleware check.
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// In a real app, you would validate the token here.
		h(w, r, ps)
	}
}

func NewRouter(userController controller.UserController, tokoController controller.TokoController) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userController.ShowLoginForm(w, r, nil)
	})
	router.POST("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userController.Login(w, r, nil)
	})

	router.GET("/dashboard", protected(userController.Dashboard))
	router.GET("/profile", protected(userController.ShowProfile))

	router.GET("/user", protected(userController.FindAll))
	router.GET("/user/add", protected(userController.Create))
	router.POST("/user/add", protected(userController.Create))
	router.GET("/user/edit/:userID", protected(userController.FindById))
	router.POST("/user/edit/:userID", protected(userController.Update))
	router.GET("/user/delete/:userID", protected(userController.Delete))

	router.GET("/toko", protected(tokoController.FindAll))
	router.GET("/toko/add", protected(tokoController.Create))
	router.POST("/toko/add", protected(tokoController.Create))
	router.GET("/toko/edit/:tokoID", protected(tokoController.FindById))
	router.POST("/toko/edit/:tokoID", protected(tokoController.Update))
	router.GET("/toko/delete/:tokoID", protected(tokoController.Delete))

	return router
}
