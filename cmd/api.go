package main

import (
	"context"
	"log"

	"os"

	"github.com/berzz26/foundry_api/internal/auth"
	"github.com/berzz26/foundry_api/internal/companies"
	"github.com/berzz26/foundry_api/internal/founders"
	"github.com/berzz26/foundry_api/internal/jobs"
	"github.com/berzz26/foundry_api/internal/outreach"
	"github.com/berzz26/foundry_api/internal/savedjobs"
	"github.com/berzz26/foundry_api/internal/users"
	"github.com/berzz26/foundry_api/pkg/config"
	"github.com/berzz26/foundry_api/pkg/database"
	"github.com/berzz26/foundry_api/pkg/email"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	cfg := config.LoadConfig()
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := auth.LoadRolePrivileges(ctx, db.DB); err != nil {
		log.Fatalf("failed to load role privileges: %v", err)
	}

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

	outreachRepo := outreach.NewRepository(db.DB)
	outreachService := outreach.NewService(outreachRepo, email.NewSender(email.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Pass:     cfg.SMTPPass,
		From:     cfg.SMTPFrom,
		FromName: cfg.SMTPFromName,
	}))
	outreachHandler := outreach.NewHandler(outreachService)

	authRepo := auth.NewRepository(db.DB)
	authService := auth.NewService(userService, authRepo)
	authHandler := auth.NewHandler(authService, cfg)

	savedJobsRepo := savedjobs.NewRepository(db.DB)
	savedJobsService := savedjobs.NewService(savedJobsRepo)
	savedJobsHandler := savedjobs.NewHandler(savedJobsService)

	app := fiber.New()
	app.Use(recover.New())

	if os.Getenv("DEV") != "" || os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "dev" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Apply optional auth parsing globally to all v1 endpoints
	v1.Use(auth.OptionalAuth())

	//mount the routes
	v1.Mount("/auth", authHandler.SetupRoutes())
	// Restrict all user route group endpoints to authenticated users
	usersGroup := v1.Group("/users", auth.RequireAuth())
	usersGroup.Mount("/", userHandler.SetupRoutes())

	v1.Mount("/companies", companyHandler.SetupRoutes())
	v1.Mount("/jobs", jobHandler.SetupRoutes())

	// Restrict all founder route group endpoints to users with founders:access privilege
	foundersGroup := v1.Group("/founders", auth.RequireAuth(), auth.RequirePrivilege("founders:access"))
	foundersGroup.Mount("/", founderHandler.SetupRoutes())

	// Outreach swipe deck + email send. Exposes founder contact info, so it is
	// restricted to users with the founders:access privilege.
	outreachGroup := v1.Group("/outreach", auth.RequireAuth(), auth.RequirePrivilege("founders:access"))
	outreachGroup.Mount("/", outreachHandler.SetupRoutes())

	// Restrict all saved-jobs endpoints to authenticated users
	savedJobsGroup := v1.Group("/saved-jobs", auth.RequireAuth())
	savedJobsGroup.Mount("/", savedJobsHandler.SetupRoutes())

	log.Fatal(app.Listen(":" + cfg.HTTPPort))

}
