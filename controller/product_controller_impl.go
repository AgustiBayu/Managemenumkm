package controller

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"Managemenumkm/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

type ProductControllerImpl struct {
	ProductService         service.ProductService
	ProductBatchService    service.ProductBatchService
	ProductCategoryService service.ProductCategoryService
}

func NewProductController(productService service.ProductService, batchService service.ProductBatchService, categoryService service.ProductCategoryService) ProductController {
	return &ProductControllerImpl{
		ProductService:         productService,
		ProductBatchService:    batchService,
		ProductCategoryService: categoryService,
	}
}

// --- Master Product Handlers ---

func (c *ProductControllerImpl) CreateView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	categories, _ := c.ProductCategoryService.FindAll(r.Context())
	data := map[string]interface{}{
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product/product_form_add.html", data)
}

func (c *ProductControllerImpl) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	price, _ := strconv.Atoi(r.PostFormValue("price"))
	catID, _ := strconv.Atoi(r.PostFormValue("category_id"))

	req := domain.ProductCreateRequest{
		Name:       r.PostFormValue("name"),
		SKU:        r.PostFormValue("sku"),
		Price:      uint(price),
		CategoryID: uint(catID),
	}

	_, err := c.ProductService.Create(req)
	if err != nil {
		// TODO: better error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product", http.StatusSeeOther)
}

func (c *ProductControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	products, _ := c.ProductService.FindAll()
	helper.RenderTemplate(w, "template/product/product_list.html", map[string]interface{}{"Products": products})
}

func (c *ProductControllerImpl) FindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("productId"))

	product, err := c.ProductService.FindById(uint(id))
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	categories, _ := c.ProductCategoryService.FindAll(r.Context())
	data := map[string]interface{}{
		"Product":    product,
		"Categories": categories,
	}
	helper.RenderTemplate(w, "template/product/product_form_edit.html", data)
}

func (c *ProductControllerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	id, _ := strconv.Atoi(ps.ByName("productId"))
	price, _ := strconv.Atoi(r.PostFormValue("price"))
	catID, _ := strconv.Atoi(r.PostFormValue("category_id"))

	req := domain.ProductUpdateRequest{
		ID:         uint(id),
		Name:       r.PostFormValue("name"),
		SKU:        r.PostFormValue("sku"),
		Price:      uint(price),
		CategoryID: uint(catID),
	}

	_, err := c.ProductService.Update(req)
	if err != nil {
		// TODO: better error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product", http.StatusSeeOther)
}

func (c *ProductControllerImpl) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("productId"))
	c.ProductService.Delete(uint(id))
	http.Redirect(w, r, "/product", http.StatusSeeOther)
}

// --- Stock & Batch Handlers ---

func (c *ProductControllerImpl) AddStockView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := map[string]interface{}{
		"Step":    1, // Initial step: scan barcode
		"Barcode": "",
	}
	helper.RenderTemplate(w, "template/product/product_stock_add.html", data)
}

func (c *ProductControllerImpl) AddStockScanBarcode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	barcode := r.PostFormValue("barcode")

	_, err := c.ProductBatchService.FindByBarcode(barcode)
	if err == nil {
		// Barcode exists, which is an error for adding a NEW batch
		data := map[string]interface{}{
			"Step":  1,
			"Error": "Barcode '" + barcode + "' sudah terdaftar di batch lain.",
		}
		helper.RenderTemplate(w, "template/product/product_stock_add.html", data)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A different database error occurred
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Barcode is new and valid, proceed to step 2
	products, _ := c.ProductService.FindAll()
	data := map[string]interface{}{
		"Step":     2, // Second step: associate product and add details
		"Barcode":  barcode,
		"Products": products,
	}
	helper.RenderTemplate(w, "template/product/product_stock_add.html", data)
}

func (c *ProductControllerImpl) AddStockSaveBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	productID, _ := strconv.Atoi(r.PostFormValue("product_id"))
	stock, _ := strconv.Atoi(r.PostFormValue("stock"))

	req := domain.AddStockStep2Request{
		Barcode:   r.PostFormValue("barcode"),
		ProductID: uint(productID),
		Exp:       r.PostFormValue("exp"),
		Stock:     uint(stock),
	}

	_, err := c.ProductService.AddStock(req)
	if err != nil {
		// Re-render form with error
		products, _ := c.ProductService.FindAll()
		data := map[string]interface{}{
			"Step":     2,
			"Barcode":  req.Barcode,
			"Products": products,
			"Error":    err.Error(),
		}
		helper.RenderTemplate(w, "template/product/product_stock_add.html", data)
		return
	}

	http.Redirect(w, r, "/product", http.StatusSeeOther)
}

func (c *ProductControllerImpl) ListBatches(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("productId"))

	product, err := c.ProductService.FindById(uint(id))
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	batches, _ := c.ProductBatchService.FindByProductID(uint(id))
	data := map[string]interface{}{
		"Product": product,
		"Batches": batches,
	}
	helper.RenderTemplate(w, "template/product/product_batch_list.html", data)
}

// --- New Batch Edit Handlers ---

func (c *ProductControllerImpl) EditBatchView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	batchID, _ := strconv.Atoi(ps.ByName("batchId"))

	batch, err := c.ProductBatchService.FindById(uint(batchID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Product batch not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	product, err := c.ProductService.FindById(batch.ProductID)
	if err != nil {
		http.Error(w, "Associated product not found", http.StatusNotFound)
		return
	}

	// Format date for input[type=date]
	// batch.Exp = helper.FormatDateForHTML(batch.Exp)

	data := map[string]interface{}{
		"Batch":   batch,
		"Product": product,
	}
	helper.RenderTemplate(w, "template/product/product_batch_form_edit.html", data)
}

func (c *ProductControllerImpl) EditBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	batchID, _ := strconv.Atoi(ps.ByName("batchId"))
	productID, _ := strconv.Atoi(r.PostFormValue("product_id"))
	stock, _ := strconv.Atoi(r.PostFormValue("stock"))

	req := domain.ProductBatchUpdateRequest{
		ID:        uint(batchID),
		ProductID: uint(productID),
		Stock:     uint(stock),
		Exp:       r.PostFormValue("exp"),
	}

	_, err := c.ProductBatchService.Update(req)
	if err != nil {
		// Re-render form with error
		batch, _ := c.ProductBatchService.FindById(uint(batchID))
		product, _ := c.ProductService.FindById(uint(productID))
		// batch.Exp = helper.FormatDateForHTML(batch.Exp)

		data := map[string]interface{}{
			"Batch":   batch,
			"Product": product,
			"Error":   err.Error(),
		}
		helper.RenderTemplate(w, "template/product/product_batch_form_edit.html", data)
		return
	}

	http.Redirect(w, r, "/product/batches/"+strconv.Itoa(productID), http.StatusSeeOther)
}