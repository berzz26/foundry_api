package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	ctx := context.Background()
	addedUser, err := h.service.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return c.JSON(addedUser)
}
