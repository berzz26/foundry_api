package founders

import (
	"context"
	"errors"
	"log"
	"net/http"
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

func (h *Handler) Create(c *fiber.Ctx) error {
	dto := new(CreateFounderDTO)
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

	var founderID int64
	if dto.ID != nil {
		founderID = *dto.ID
	}

	founder := &Founder{
		ID:           founderID,
		CompanyID:    dto.CompanyID,
		FullName:     dto.FullName,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Bio:          dto.Bio,
		Linkedin:     dto.Linkedin,
		Twitter:      dto.Twitter,
		AvatarURL:    dto.AvatarURL,
		AvatarThumb:  dto.AvatarThumb,
		AvatarMedium: dto.AvatarMedium,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	created, err := h.service.Create(ctx, founder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create founder",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mapToResponseDTO(created))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing founder ID parameter",
		})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid founder ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	founder, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, ErrFounderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Founder not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(mapToResponseDTO(founder))
}
func (h *Handler) GetByCompanyID(c *fiber.Ctx) error {
	companyIDStr := c.Params("companyId")
	if companyIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing company ID parameter",
		})
	}

	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	founders, err := h.service.GetByCompanyID(ctx, companyID)
	if err != nil {
		if errors.Is(err, ErrFounderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Founder not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(founders)
}
func (h *Handler) List(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	companyIDStr := c.Query("companyId")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	if limit > 100 {
		limit = 100
	}

	var companyID *int64
	if companyIDStr != "" {
		if cid, err := strconv.ParseInt(companyIDStr, 10, 64); err == nil {
			companyID = &cid
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid company ID format",
			})
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	foundersList, err := h.service.List(ctx, companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list founders",
		})
	}

	response := make([]FounderResponseDTO, len(foundersList))
	for i, f := range foundersList {
		response[i] = mapToResponseDTO(&f)
	}

	return c.JSON(response)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing founder ID parameter",
		})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid founder ID format",
		})
	}

	dto := new(UpdateFounderDTO)
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

	existing, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrFounderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Founder not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if dto.CompanyID != nil {
		existing.CompanyID = dto.CompanyID
	}
	if dto.FullName != nil {
		existing.FullName = *dto.FullName
	}
	if dto.FirstName != nil {
		existing.FirstName = dto.FirstName
	}
	if dto.LastName != nil {
		existing.LastName = dto.LastName
	}
	if dto.Bio != nil {
		existing.Bio = dto.Bio
	}
	if dto.Linkedin != nil {
		existing.Linkedin = dto.Linkedin
	}
	if dto.Twitter != nil {
		existing.Twitter = dto.Twitter
	}
	if dto.AvatarURL != nil {
		existing.AvatarURL = dto.AvatarURL
	}
	if dto.AvatarThumb != nil {
		existing.AvatarThumb = dto.AvatarThumb
	}
	if dto.AvatarMedium != nil {
		existing.AvatarMedium = dto.AvatarMedium
	}

	updated, err := h.service.Update(ctx, existing)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update founder",
		})
	}

	return c.JSON(mapToResponseDTO(updated))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing founder ID parameter",
		})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid founder ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if _, err := h.service.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrFounderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Founder not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete founder",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

func mapToResponseDTO(f *Founder) FounderResponseDTO {
	return FounderResponseDTO{
		ID:           f.ID,
		CompanyID:    f.CompanyID,
		FullName:     f.FullName,
		FirstName:    f.FirstName,
		LastName:     f.LastName,
		Bio:          f.Bio,
		Linkedin:     f.Linkedin,
		Twitter:      f.Twitter,
		AvatarURL:    f.AvatarURL,
		AvatarSourceURL: f.AvatarSourceURL,
		AvatarThumb:  f.AvatarThumb,
		AvatarSourceThumb: f.AvatarSourceThumb,
		AvatarMedium: f.AvatarMedium,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
}
