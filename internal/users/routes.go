package users

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(
	router fiber.Router,
	handler *Handler,
) {
	users := router.Group("/users")

	users.Post("/", handler.AddUser)
	// users.Get("/me", handler.GetMe)

}
