package savedjobs

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/", h.Save)
	app.Get("/", h.List)
	app.Get("/:jobId", h.IsSaved)
	app.Delete("/:jobId", h.Delete)

	return app
}
