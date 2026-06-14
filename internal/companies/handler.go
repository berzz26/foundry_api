package companies

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	filters := CompanyFilters{}
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Default pagination values if not provided
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	response, err := h.service.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve companies",
		})
	}

	return c.JSON(response)
}

func (h *Handler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing parameter Slug",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	company, err := h.service.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(mapToDetailResponse(company))
}

func (h *Handler) GetMeta(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	meta, err := h.service.GetMetadata(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve metadata",
		})
	}

	return c.JSON(meta)
}

func mapToDetailResponse(c *Company) CompanyDetailResponse {
	return CompanyDetailResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Slug:               c.Slug,
		Website:            c.Website,
		Tagline:            c.Tagline,
		Description:        c.Description,
		HiringDescription:  c.HiringDescription,
		TechStack:          c.TechStack,
		Batch:              c.Batch,
		Stage:              c.Stage,
		TeamSize:           c.TeamSize,
		Location:           c.Location,
		ParentSector:       c.ParentSector,
		ChildSector:        c.ChildSector,
		Industry:           c.Industry,
		LogoURL:            c.LogoURL,
		SourceLogoURL:      c.SourceLogoURL,
		SmallLogoURL:       c.SmallLogoURL,
		SourceSmallLogoURL: c.SourceSmallLogoURL,
		Country:            c.Country,
		FoundedAt:          c.FoundedAt,
		LinkedinURL:        c.LinkedinURL,
		TwitterURL:         c.TwitterURL,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
