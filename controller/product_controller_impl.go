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

type ProductControllerImpl struct {
	ProductService         service.ProductService
	ProductCategoryService service.ProductCategoryService
}

func NewProductController(productService service.ProductService, productCategoryService service.ProductCategoryService) ProductController {
	return &ProductControllerImpl{
		ProductService:         productService,
		ProductCategoryService: productCategoryService,
	}
}

func (c *ProductControllerImpl) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		priceStr := r.FormValue("price")
		exp := r.FormValue("exp")
		stockStr := r.FormValue("stock")
		categoryIdStr := r.FormValue("category_id")
		barcode := r.FormValue("barcode")

		price, err := strconv.Atoi(priceStr)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "Invalid category ID format", http.StatusBadRequest)
			return
		}
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "Invalid stock format", http.StatusBadRequest)
			return
		}
		req := &domain.ProductCreateRequest{Name: name,
			Price: uint(price), Exp: exp, Stock: uint(stock), CategoryID: uint(categoryId), Barcode: barcode}

		file, handler, err := r.FormFile("thumbnail")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Failed to get thumbnail file", http.StatusBadRequest)
			return
		}
		if file != nil {
			defer file.Close()
		}

		if err := c.ProductService.Create(context.Background(), req, file, handler); err != nil {
			http.Error(w, "Gagal menyimpan data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/product", http.StatusSeeOther)
		return
	}
	categories, _ := c.ProductCategoryService.FindAll(context.Background())

	// Siapkan data untuk template
	data := map[string]interface{}{
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product/product_form_add.html", data)
}

func (c *ProductControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	products, _ := c.ProductService.FindAll(context.Background())

	helper.RenderTemplate(w, "template/product/product_list.html", map[string]interface{}{"Products": products})
}

func (c *ProductControllerImpl) FindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("productId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := c.ProductService.FindById(context.Background(), id)
	if err != nil {
		http.Error(w, "Failed to find product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := c.ProductCategoryService.FindAll(context.Background())
	if err != nil {
		http.Error(w, "Failed to find categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Product":    product,
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product/product_form_edit.html", data)
}

func (c *ProductControllerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == http.MethodPost {
		idStr := ps.ByName("productId")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		barcode := r.FormValue("barcode")
		priceStr := r.FormValue("price")
		exp := r.FormValue("exp")
		stockStr := r.FormValue("stock")
		categoryIdStr := r.FormValue("category_id")

		price, err := strconv.Atoi(priceStr)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "Invalid category ID format", http.StatusBadRequest)
			return
		}
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "Invalid stock format", http.StatusBadRequest)
			return
		}

		req := &domain.ProductUpdateRequest{
			ID:         uint(id),
			Name:       name,
			Barcode:    barcode,
			Price:      uint(price),
			Exp:        exp,
			Stock:      uint(stock),
			CategoryID: uint(categoryId),
		}

		file, handler, err := r.FormFile("thumbnail")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Failed to get thumbnail file", http.StatusBadRequest)
			return
		}
		if file != nil {
			defer file.Close()
		}

		if err := c.ProductService.Update(context.Background(), req, file, handler); err != nil {
			http.Error(w, "Gagal memperbarui data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/product", http.StatusSeeOther)
		return
	}

	idStr := ps.ByName("productId")
	id, _ := strconv.Atoi(idStr)

	product, _ := c.ProductService.FindById(context.Background(), id)
	categories, _ := c.ProductCategoryService.FindAll(context.Background())
	data := map[string]interface{}{
		"Product":    product,
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product/product_form_edit.html", data)
}

func (c *ProductControllerImpl) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("productId")
	id, _ := strconv.Atoi(idStr)

	c.ProductService.Delete(context.Background(), id)

	http.Redirect(w, r, "/product", http.StatusSeeOther)
}
