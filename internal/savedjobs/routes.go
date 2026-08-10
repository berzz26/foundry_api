package savedjobs

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/", h.SaveJob)
	app.Get("/", h.List)
	app.Post("/companies", h.SaveCompany)
	app.Get("/companies/:companyId", h.IsSavedCompany)
	app.Delete("/companies/:companyId", h.DeleteCompany)
	app.Get("/:jobId", h.IsSavedJob)
	app.Delete("/:jobId", h.DeleteJob)

	return app
}
