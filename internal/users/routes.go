package users

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/", h.AddUser)
	app.Get("/", h.List)
	app.Get("/:id", h.GetByID)
	app.Get("/email/:email", h.GetByEmail)
	app.Put("/:id", h.Update)
	app.Delete("/:id", h.Delete)

	return app
}
