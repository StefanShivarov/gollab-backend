package common

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func HealthRoute(r chi.Router, db *gorm.DB) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		db, err := db.DB()
		if err != nil || db.Ping() != nil {
			WriteError(w, &ApiError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Database not available!",
			})
			return
		}
		WriteJSON(w, http.StatusOK, []byte(`{"status":"ok"}`))
	})
}
