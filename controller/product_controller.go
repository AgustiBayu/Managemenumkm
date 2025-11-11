package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProductController interface {
	// Master Product handlers
	CreateView(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

	// Stock & Batch handlers
	AddStockView(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	AddStockScanBarcode(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	AddStockSaveBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ListBatches(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

	// New Batch Edit Handlers
	EditBatchView(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	EditBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}