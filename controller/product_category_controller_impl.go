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

type ProductCategoryControllerImpl struct {
	ProductCategoryService service.ProductCategoryService
}

func NewProductCategoryController(productCategoryService service.ProductCategoryService) ProductCategoryController {
	return &ProductCategoryControllerImpl{
		ProductCategoryService: productCategoryService,
	}
}

func (c *ProductCategoryControllerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == http.MethodPost {
		category := r.FormValue("category")
		req := &domain.ProductCategoryCreateRequest{Category: category}
		_ = c.ProductCategoryService.Create(context.Background(), req)

		http.Redirect(w, r, "/category", http.StatusSeeOther)
		return
	}
	helper.RenderTemplate(w, "template/product_category/produk_kategori_add.html", nil)
}

func (c *ProductCategoryControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	categories, _ := c.ProductCategoryService.FindAll(context.Background())
	data := map[string]interface{}{
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product_category/produk_kategori_list.html", data)
}

func (c *ProductCategoryControllerImpl) FindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("categoryID"))
	category, _ := c.ProductCategoryService.FindById(context.Background(), id)

	data := map[string]interface{}{
		"Category": category,
	}
	helper.RenderTemplate(w, "template/product_category/produk_kategori_update.html", data)
}

func (c *ProductCategoryControllerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("categoryID"))

	if r.Method == http.MethodPost {
		categoryValue := r.FormValue("category")
		req := &domain.ProductCategoryUpdateRequest{
			ID:       uint(id),
			Category: categoryValue,
		}
		if err := c.ProductCategoryService.Update(context.Background(), req); err != nil {
			http.Error(w, "Gagal memperbarui data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/category", http.StatusSeeOther)
		return
	}

	category, err := c.ProductCategoryService.FindById(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Category": category,
	}
	helper.RenderTemplate(w, "template/product_category/produk_kategori_update.html", data)
}

func (c *ProductCategoryControllerImpl) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("categoryID"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err = c.ProductCategoryService.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/category", http.StatusSeeOther)
}
