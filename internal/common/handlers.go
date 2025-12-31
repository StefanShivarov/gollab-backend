package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error) {
	var apiError *ApiError
	w.Header().Set("Content-Type", "application/json")

	if errors.As(err, &apiError) {
		w.WriteHeader(apiError.StatusCode)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{Error: *apiError})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(&ErrorResponse{
		Error: *InternalServerError("Unexpected error occurred!"),
	})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
