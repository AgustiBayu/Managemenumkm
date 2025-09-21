package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/service"
	"context"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TokoControllerImpl struct {
	TokoService service.TokoService
}

func NewTokoController(tokoService service.TokoService) TokoController {
	return &TokoControllerImpl{TokoService: tokoService}
}

func (t *TokoControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		address := r.FormValue("address")
		req := &domain.TokoCreateRequest{Name: name, Address: address}

		if err := t.TokoService.Create(context.Background(), req); err != nil {
			http.Error(w, "Gagal menyimpan data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/toko", http.StatusSeeOther)
		return
	}
	helper.RenderTemplate(w, "template/toko/toko_add.html", nil)
}

func (t *TokoControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 10
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}
	tokos, totalItem, err := t.TokoService.FindAll(context.Background(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Tokos":     tokos,
		"TotalItem": totalItem,
		"Page":      page,
		"PageSize":  pageSize,
	}

	helper.RenderTemplate(w, "template/toko/toko_list.html", data)
}

func (t *TokoControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("tokoID"))
	toko, _ := t.TokoService.FindById(context.Background(), id)

	data := map[string]interface{}{
		"Toko": toko,
	}
	helper.RenderTemplate(w, "template/toko/toko_update.html", data)
}

func (t *TokoControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("tokoID"))

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		address := r.FormValue("address")
		req := &domain.TokoUpdateRequest{
			ID:      uint(id),
			Name:    name,
			Address: address,
		}

		if err := t.TokoService.Update(context.Background(), req); err != nil {
			http.Error(w, "Gagal memperbarui data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/toko", http.StatusSeeOther)
		return
	}
	toko, _ := t.TokoService.FindById(context.Background(), id)
	data := map[string]interface{}{
		"Toko": toko,
	}
	helper.RenderTemplate(w, "template/toko/toko_update.html", data)
}

func (t *TokoControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idStr := params.ByName("tokoID")
	id, _ := strconv.Atoi(idStr)

	t.TokoService.Delete(context.Background(), id)

	http.Redirect(w, r, "/toko", http.StatusSeeOther)
}
