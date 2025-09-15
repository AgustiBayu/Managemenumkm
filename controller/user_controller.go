package controller

import "net/http"

type UserController interface {
	ShowLoginForm(w http.ResponseWriter, r *http.Request, params map[string]string)
	Login(w http.ResponseWriter, r *http.Request, params map[string]string)
	Dashboard(w http.ResponseWriter, r *http.Request, params map[string]string)
}
