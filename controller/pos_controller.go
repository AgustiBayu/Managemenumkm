package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PosController interface {
	ShowPosPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}