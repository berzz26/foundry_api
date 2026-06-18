package auth

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/signup", h.Signup)
	app.Post("/login", h.Login)
	app.Post("/logout", h.Logout)
	app.Get("/me", RequireAuth(), h.GetMe)

	app.Get("/google/login", h.GoogleLogin)
	app.Get("/google/callback", h.GoogleCallback)

	app.Get("/github/login", h.GithubLogin)
	app.Get("/github/callback", h.GithubCallback)

	return app	
}
