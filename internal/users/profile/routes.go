package profile

import (
	"github.com/berzz26/foundry_api/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/onboarding", auth.RequireAuth(), h.CompleteOnboarding)
	app.Get("/profile", auth.RequireAuth(), h.GetProfile)

	return app
}
