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

func (h *Handler) SaveJob(c *fiber.Ctx) error {
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

	res, err := h.service.SaveJob(ctx, userID, dto.JobID)
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

func (h *Handler) SaveCompany(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	dto := new(SaveCompanyRequest)
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

	res, err := h.service.SaveCompany(ctx, userID, dto.CompanyID)
	if err != nil {
		if errors.Is(err, ErrSavedCompanyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save company",
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
			"error": "Failed to retrieve saved items",
		})
	}

	return c.JSON(res)
}

func (h *Handler) IsSavedJob(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	jobID, err := parseID(c, "jobId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	saved, err := h.service.IsJobSaved(ctx, userID, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check saved status",
		})
	}

	return c.JSON(SavedStatusResponse{Saved: saved})
}

func (h *Handler) IsSavedCompany(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	companyID, err := parseID(c, "companyId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	saved, err := h.service.IsCompanySaved(ctx, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check saved status",
		})
	}

	return c.JSON(SavedStatusResponse{Saved: saved})
}

func (h *Handler) DeleteJob(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	jobID, err := parseID(c, "jobId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteJob(ctx, userID, jobID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove saved job",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *Handler) DeleteCompany(c *fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	companyID, err := parseID(c, "companyId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteCompany(ctx, userID, companyID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove saved company",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func parseID(c *fiber.Ctx, param string) (int64, error) {
	return strconv.ParseInt(c.Params(param), 10, 64)
}
