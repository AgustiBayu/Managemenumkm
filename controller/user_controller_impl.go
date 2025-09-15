package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/service"
	"net/http"
	"text/template"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{UserService: userService}
}

func (c *UserControllerImpl) ShowLoginForm(w http.ResponseWriter, r *http.Request, params map[string]string) {
	tmpl, err := template.ParseFiles("template/auth/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (u *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, params map[string]string) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	request := domain.LoginRequest{
		Email:    email,
		Password: password,
	}

	token, err := u.UserService.Login(&request)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (c *UserControllerImpl) Dashboard(w http.ResponseWriter, r *http.Request, params map[string]string) {
	tmpl := template.Must(template.ParseFiles("template/dashboard/dashboard.html"))
	tmpl.Execute(w, nil)
}
