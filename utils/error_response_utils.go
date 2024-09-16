package utils

import (
	"encoding/json"
	"net/http"
	"tender-service/models"
)

func WriteErrorResponse(w http.ResponseWriter, reason string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Reason: reason})
}