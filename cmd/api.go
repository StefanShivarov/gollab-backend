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
	Config    config.Config
	DB        *gorm.DB
	Validator *validator.Validate
}

func NewApplication(cfg config.Config) (*Application, error) {
	gormDB, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := db.PerformMigration(gormDB); err != nil {
		return nil, err
	}

	return &Application{
		Config:    cfg,
		DB:        gormDB,
		Validator: validator.New(),
	}, nil
}

func (app *Application) mountRoutes(r chi.Router) {
	userService := org.NewUserService(org.NewUserRepository(app.DB), app.Validator)
	userHandler := org.NewUserHandler(userService)
	teamHandler := org.NewTeamHandler(org.NewTeamService(org.NewTeamRepository(app.DB), userService, app.Validator))

	org.UserRoutes(r, userHandler)
	org.TeamRoutes(r, teamHandler)
}

func (app *Application) Routes() http.Handler {
	router := chi.NewRouter()
	app.mountRoutes(router)
	return router
}

func (app *Application) Run(addr string) {
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Routes()))
}
