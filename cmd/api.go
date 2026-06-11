package main

import (
	"log"

	"github.com/berzz26/foundry_api/internal/users"
	"github.com/berzz26/foundry_api/pkg/config"
	"github.com/berzz26/foundry_api/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.LoadConfig()
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := users.NewRepository(db.DB)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Mount("/users", userHandler.SetupRoutes())

	log.Fatal(app.Listen(":" + cfg.HTTPPort))

}
