package outreach

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	var filters ListFilters
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	res, err := h.service.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve outreach cards",
		})
	}

	return c.JSON(res)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := parseOutreachID(c)
	if err != nil {
		return err
	}

	var founderID *int64
	if v := c.Query("founderId"); v != "" {
		fid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid founderId query parameter",
			})
		}
		founderID = &fid
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	card, err := h.service.GetByID(ctx, id, founderID)
	if err != nil {
		if errors.Is(err, ErrOutreachNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Outreach not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve outreach card",
		})
	}

	return c.JSON(card)
}

func (h *Handler) Send(c *fiber.Ctx) error {
	var outreachID int64
	if idStr := c.Params("id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid outreach ID format",
			})
		}
		outreachID = id
	}

	dto := new(SendEmailRequest)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		userID = "anonymous"
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	res, err := h.service.Send(ctx, outreachID, *dto, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOutreachNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Outreach not found",
			})
		case errors.Is(err, ErrNoRecipient), errors.Is(err, ErrNoFounderEmail):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "Failed to send email: " + err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func parseOutreachID(c *fiber.Ctx) (int64, error) {
	idStr := c.Params("id")
	if idStr == "" {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing outreach ID parameter",
		})
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid outreach ID format",
		})
	}
	return id, nil
}
