package companies

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Get("/", h.List)
	app.Get("/meta", h.GetMeta)
	app.Get("/:idOrSlug", h.GetByIDOrSlug)

	return app
}
