package jobs

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SetupRoutes() *fiber.App {
	router := fiber.New()

	router.Get("/", h.List)
	router.Get("/featured", h.GetFeatured)
	router.Get("/:id", h.GetByID)
	router.Get("/:id/related", h.GetRelated)

	return router
}
