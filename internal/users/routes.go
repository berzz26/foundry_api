package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(
	router fiber.Router,
	handler *Handler,
) {
	users := router.Group("/users")

	users.Post("/", handler.AddUser)
	users.Get("/", handler.List)
	users.Get("/:id", handler.GetByID)
	users.Get("/email/:email", handler.GetByEmail)
	users.Put("/:id", handler.Update)
	users.Delete("/:id", handler.Delete)
}
