package jobs

import (
	"context"
	"strconv"
	"time"

	"github.com/berzz26/foundry_api/internal/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}



func (h *Handler) List(c *fiber.Ctx) error {
	hasList50Privilege := auth.HasPrivilege(c, "jobs:list_50")

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if !hasList50Privilege {
		res, err := h.service.GetRandomJobs(ctx, 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve random jobs",
			})
		}
		return c.JSON(res)
	}

	var filters JobFilters
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}

	maxAllowed := 50
	offset := (filters.Page - 1) * filters.Limit
	if offset >= maxAllowed {
		return c.JSON(JobListResponse{
			Jobs: []JobCardResponse{},
			Pagination: PaginationResponse{
				Page:    filters.Page,
				Limit:   filters.Limit,
				Total:   maxAllowed,
				HasNext: false,
			},
		})
	}

	if offset+filters.Limit > maxAllowed {
		filters.Limit = maxAllowed - offset
	}

	res, err := h.service.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve jobs",
		})
	}

	if res.Pagination.Total > maxAllowed {
		res.Pagination.Total = maxAllowed
	}
	res.Pagination.HasNext = (res.Pagination.Page * res.Pagination.Limit) < res.Pagination.Total

	if len(res.Jobs) > res.Pagination.Limit {
		res.Jobs = res.Jobs[:res.Pagination.Limit]
	}

	return c.JSON(res)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	job, err := h.service.GetByID(ctx, id)
	if err != nil {
		// Handle pgx.ErrNoRows or standard error
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Job not found",
		})
	}

	return c.JSON(job)
}

func (h *Handler) GetRelated(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.service.GetRelated(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve related jobs",
		})
	}

	return c.JSON(res)
}

func (h *Handler) GetFeatured(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.service.GetFeatured(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve featured jobs",
		})
	}

	return c.JSON(res)
}
