package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/service"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
	TokoService service.TokoService
}

func NewUserController(userService service.UserService, tokoService service.TokoService) UserController {
	return &UserControllerImpl{UserService: userService, TokoService: tokoService}
}

func (u *UserControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodPost {
		fname := r.FormValue("first_name")
		lname := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password") // Read password
		alamat := r.FormValue("alamat")
		number := r.FormValue("number")
		role := r.FormValue("role")
		tokoID, _ := strconv.Atoi(r.FormValue("toko_id"))
		subscriptionStatus := r.FormValue("subscription_status")
		subscriptionExpiry := r.FormValue("subscription_expiry")
		req := &domain.UserCreateRequest{FName: fname, LName: lname, Email: email, Password: password, Alamat: alamat, Number: number,
			Role: domain.Role(role), TokoID: uint(tokoID), SubscriptionStatus: subscriptionStatus, SubscriptionExpiry: subscriptionExpiry}

		if err := u.UserService.Create(context.Background(), req); err != nil {
			http.Error(w, "Gagal menyimpan data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	tokos, _, _ := u.TokoService.FindAll(context.Background(), 1, 100) // Fetch all tokos
	data := map[string]interface{}{
		"Tokos": tokos,
	}
	helper.RenderTemplate(w, "template/user/user_add.html", data)
}

func (u *UserControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 100
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}
	users, totalItem, err := u.UserService.FindAll(context.Background(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Users":     users,
		"TotalItem": totalItem,
		"Page":      page,
		"PageSize":  pageSize,
	}

	helper.RenderTemplate(w, "template/user/user_list.html", data)
}
func (u *UserControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("userID"))
	user, _ := u.UserService.FindById(context.Background(), id)
	tokos, _, _ := u.TokoService.FindAll(context.Background(), 1, 100) // Fetch all tokos

	data := map[string]interface{}{
		"User":  user,
		"Tokos": tokos,
	}
	helper.RenderTemplate(w, "template/user/user_update.html", data)
}
func (u *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("userID"))

	if r.Method == http.MethodPost {
		fname := r.FormValue("first_name")
		lname := r.FormValue("last_name")
		email := r.FormValue("email")
		alamat := r.FormValue("alamat")
		number := r.FormValue("number")
		role := r.FormValue("role")
		tokoID, _ := strconv.Atoi(r.FormValue("toko_id"))
		password := r.FormValue("password") // Read password
		subscriptionStatus := r.FormValue("subscription_status")
		subscriptionExpiry := r.FormValue("subscription_expiry")
		req := &domain.UserUpdateRequest{
			ID:                 uint(id),
			FName:              fname,
			LName:              lname,
			Email:              email,
			Password:           password, // Add password to request
			Alamat:             alamat,
			Number:             number,
			Role:               domain.Role(role),
			TokoID:             uint(tokoID),
			SubscriptionStatus: subscriptionStatus,
			SubscriptionExpiry: subscriptionExpiry,
		}

		file, handler, err := r.FormFile("thumbnail")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Failed to get thumbnail file", http.StatusBadRequest)
			return
		}
		if file != nil {
			defer file.Close()
		}
		if err := u.UserService.Update(context.Background(), req, file, handler); err != nil {
			http.Error(w, "Gagal memperbarui data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}
	user, _ := u.UserService.FindById(context.Background(), id)
	tokos, _, _ := u.TokoService.FindAll(context.Background(), 1, 100)
	data := map[string]interface{}{
		"User":  user,
		"Tokos": tokos,
	}
	helper.RenderTemplate(w, "template/user/user_update.html", data)
}
func (u *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idStr := params.ByName("userID")
	id, _ := strconv.Atoi(idStr)

	u.UserService.Delete(context.Background(), id)

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (u *UserControllerImpl) ShowLoginForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	helper.RenderSingleTemplate(w, "template/auth/login.html", nil)
}

func (u *UserControllerImpl) ShowProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tokenString := cookie.Value
	claims, err := helper.ValidateJWT(tokenString)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := u.UserService.FindByEmail(context.Background(), claims.Email)
	if err != nil {
		// Handle error, maybe redirect to login or show an error page
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	now := time.Now()
	totalDuration := user.SubscriptionExpiryDate.Sub(user.CreatedAt)
	remainingDuration := user.SubscriptionExpiryDate.Sub(now)

	progressPercent := 0.0
	if totalDuration > 0 {
		progressPercent = (float64(remainingDuration) / float64(totalDuration)) * 100
	}

	if progressPercent < 0 {
		progressPercent = 0
	}
	if progressPercent > 100 {
		progressPercent = 100
	}
	progressColor := "bg-primary"
	if progressPercent > 80 {
		progressColor = "bg-danger"
	} else if progressPercent > 50 {
		progressColor = "bg-warning"
	}

	data := map[string]interface{}{
		"User": user,
		"Progress": map[string]interface{}{
			"Percent": int(progressPercent),
			"Color":   progressColor,
		},
	}
	helper.RenderSingleTemplate(w, "template/profile/profile.html", data)
}

func (u *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (u *UserControllerImpl) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Mengatur waktu kedaluwarsa
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (u *UserControllerImpl) Dashboard(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	helper.RenderTemplate(w, "template/dashboard/dashboard.html", nil)
}
