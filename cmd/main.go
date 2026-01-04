package main

import (
	"fmt"

	"github.com/StefanShivarov/gollab-backend/internal/config"
)

func main() {
	cfg := config.Load()

	app, err := NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	app.Run(fmt.Sprintf(":%d", cfg.ApiPort))
}
