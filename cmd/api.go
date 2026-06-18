package main

import (
	"log"

	"os"

	"github.com/berzz26/foundry_api/internal/companies"
	"github.com/berzz26/foundry_api/internal/founders"
	"github.com/berzz26/foundry_api/internal/jobs"
	"github.com/berzz26/foundry_api/internal/users"
	"github.com/berzz26/foundry_api/pkg/config"
	"github.com/berzz26/foundry_api/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	companyRepo := companies.NewRepository(db.DB)
	companyService := companies.NewService(companyRepo)
	companyHandler := companies.NewHandler(companyService)

	jobRepo := jobs.NewRepository(db.DB)
	jobService := jobs.NewService(jobRepo)
	jobHandler := jobs.NewHandler(jobService)

	founderRepo := founders.NewRepository(db.DB)
	founderService := founders.NewService(founderRepo)
	founderHandler := founders.NewHandler(founderService)

	app := fiber.New()
	app.Use(recover.New())

	if os.Getenv("DEV") != "" || os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "dev" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//mount the routes
	v1.Mount("/users", userHandler.SetupRoutes())
	v1.Mount("/companies", companyHandler.SetupRoutes())
	v1.Mount("/jobs", jobHandler.SetupRoutes())
	v1.Mount("/founders", founderHandler.SetupRoutes())

	log.Fatal(app.Listen(":" + cfg.HTTPPort))

}
