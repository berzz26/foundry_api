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
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	signedIn := c.Locals("user_id") != nil
	maxAllowed := 10
	if signedIn {
		maxAllowed = 50
	}

	if filters.Offset >= maxAllowed {
		return c.JSON(CompanyListResponse{
			Companies: []CompanyCardResponse{},
			Pagination: PaginationResponse{
				Total:   int64(maxAllowed),
				Limit:   filters.Limit,
				Offset:  filters.Offset,
				HasNext: false,
			},
		})
	}

	if filters.Offset+filters.Limit > maxAllowed {
		filters.Limit = maxAllowed - filters.Offset
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	response, err := h.service.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve companies",
		})
	}

	// Adjust total and hasNext pagination indicators in output
	if response.Pagination.Total > int64(maxAllowed) {
		response.Pagination.Total = int64(maxAllowed)
	}
	response.Pagination.HasNext = int64(response.Pagination.Offset+response.Pagination.Limit) < response.Pagination.Total

	if len(response.Companies) > response.Pagination.Limit {
		response.Companies = response.Companies[:response.Pagination.Limit]
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
		CompanyDescription: c.Description,
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
		FacebookURL:        c.FacebookURL,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
