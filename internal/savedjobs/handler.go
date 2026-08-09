package savedjobs

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	validate        = validator.New()
	ErrUnauthorized = errors.New("unauthorized")
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func userIDFromCtx(c *fiber.Ctx) (string, error) {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}

func (h *Handler) Save(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	dto := new(SaveJobRequest)
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

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.service.Save(ctx, userID, dto.JobID)
	if err != nil {
		if errors.Is(err, ErrSavedJobNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Job not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save job",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *Handler) List(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.service.ListByUser(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve saved jobs",
		})
	}

	return c.JSON(res)
}

func (h *Handler) IsSaved(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	jobID, err := parseJobID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	saved, err := h.service.IsSaved(ctx, userID, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check saved status",
		})
	}

	return c.JSON(SavedStatusResponse{Saved: saved})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	jobID, err := parseJobID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, userID, jobID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove saved job",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func parseJobID(c *fiber.Ctx) (int64, error) {
	return strconv.ParseInt(c.Params("jobId"), 10, 64)
}
