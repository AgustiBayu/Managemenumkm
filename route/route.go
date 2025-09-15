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

func NewRouter(userController controller.UserController) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userController.ShowLoginForm(w, r, nil)
	})
	router.POST("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userController.Login(w, r, nil)
	})
	router.GET("/dashboard", protected(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userController.Dashboard(w, r, nil)
	}))
	return router
}
