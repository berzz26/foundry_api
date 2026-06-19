package profile

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CompleteOnboarding(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var req CompleteOnboardingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.service.CompleteOnboarding(ctx, userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	p, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Profile not found",
		})
	}

	// Safety check to ensure arrays are not nil (return [] instead of null in JSON)
	if p.Goals == nil {
		p.Goals = []string{}
	}
	if p.InterestedRoles == nil {
		p.InterestedRoles = []string{}
	}
	if p.PreferredLocations == nil {
		p.PreferredLocations = []string{}
	}
	if p.CompanyStagePreferences == nil {
		p.CompanyStagePreferences = []string{}
	}

	return c.JSON(UserProfileResponse{
		UserType:                p.UserType,
		ExperienceLevel:         p.ExperienceLevel,
		Goals:                   p.Goals,
		InterestedRoles:         p.InterestedRoles,
		PreferredLocations:      p.PreferredLocations,
		CompanyStagePreferences: p.CompanyStagePreferences,
	})
}
