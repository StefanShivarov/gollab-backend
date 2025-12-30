package main

import "github.com/StefanShivarov/gollab-backend/internal/config"

func main() {
	cfg := config.Load()

	app, err := NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	app.Run(":8080")
}
