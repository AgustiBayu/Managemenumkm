package helper

import (
	"encoding/json"
	"net/http"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{
		"error": message,
		"status": "error",
	}
	json.NewEncoder(w).Encode(response)
}

func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"data": data,
		"status": "success",
	}
	json.NewEncoder(w).Encode(response)
}
