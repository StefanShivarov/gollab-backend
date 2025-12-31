package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StefanShivarov/gollab-backend/internal/config"
	"github.com/StefanShivarov/gollab-backend/internal/db"
	"github.com/StefanShivarov/gollab-backend/internal/org"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Application struct {
	Config      config.Config
	DB          *gorm.DB
	UserHandler *org.UserHandler
}

func NewApplication(cfg config.Config) (*Application, error) {
	gormDB, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := db.PerformMigration(gormDB); err != nil {
		return nil, err
	}

	v := validator.New()
	userHandler := org.NewUserHandler(org.NewUserService(org.NewUserRepository(gormDB), v))

	return &Application{
		Config:      cfg,
		DB:          gormDB,
		UserHandler: userHandler,
	}, nil
}

func (app *Application) Routes() http.Handler {
	router := chi.NewRouter()
	org.Routes(router, app.UserHandler)
	return router
}

func (app *Application) Run(addr string) {
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Routes()))
}
