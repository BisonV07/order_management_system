package helpers

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse writes a JSON response with the given status code
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes an error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string, message string) {
	WriteJSONResponse(w, statusCode, map[string]string{
		"error":   errorMsg,
		"message": message,
	})
}

