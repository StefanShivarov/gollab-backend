package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StefanShivarov/gollab-backend/internal/config"
	"github.com/StefanShivarov/gollab-backend/internal/db"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Application struct {
	Config config.Config
	DB     *gorm.DB
}

func NewApplication(cfg config.Config) (*Application, error) {
	gormDB, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &Application{
		Config: cfg,
		DB:     gormDB,
	}, nil
}

func (app *Application) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		sqlDb, err := app.DB.DB()
		status := "ok"
		if err != nil || sqlDb.Ping() != nil {
			w.WriteHeader(http.StatusInternalServerError)
			status = "unhealthy"
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "` + status + `"}`))
	})

	return router
}

func (app *Application) Run(addr string) {
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Routes()))
}
