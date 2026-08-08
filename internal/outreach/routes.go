package outreach

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Get("/", h.List)
	app.Get("/:id", h.GetByID)
	app.Post("/:id/send", h.Send)

	return app
}
