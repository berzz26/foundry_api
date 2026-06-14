package companies

import (
	"context"
	"errors"
	"strconv"
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

	list, err := h.service.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve companies",
		})
	}

	response := make([]CompanyResponseDTO, len(list))
	for i, comp := range list {
		response[i] = mapToResponseDTO(&comp)
	}

	return c.JSON(response)
}

func (h *Handler) GetByIDOrSlug(c *fiber.Ctx) error {
	param := c.Params("idOrSlug")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing parameter ID or Slug",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var company *Company
	var err error

	// If it parses as int64, look up by ID first
	if id, parseErr := strconv.ParseInt(param, 10, 64); parseErr == nil {
		company, err = h.service.GetByID(ctx, id)
		if err != nil && !errors.Is(err, ErrCompanyNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	// If lookup by ID didn't find the company (or param wasn't an int), try lookup by Slug
	if company == nil {
		company, err = h.service.GetBySlug(ctx, param)
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
	}

	return c.JSON(mapToResponseDTO(company))
}

func mapToResponseDTO(c *Company) CompanyResponseDTO {
	return CompanyResponseDTO{
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
