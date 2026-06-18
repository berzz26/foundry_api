package users

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddUser(c *fiber.Ctx) error {
	dto := new(AddUserDTO)

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
	if dto.Provider == "" {
		dto.Provider = "local"
	}
	user := &User{
		Email:           dto.Email,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		PasswordHash:    dto.Password,
		ProfileImageURL: dto.ProfileImageURL,
		Provider:        dto.Provider,
		ProviderID:      dto.ProviderID,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	addedUser, err := h.service.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mapToResponseDTO(addedUser))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing user ID parameter",
		})
	}

	loggedInID := c.Locals("user_id")
	loggedInRole := c.Locals("user_role")
	if loggedInID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if loggedInRole != "admin" && loggedInID.(string) != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden. Insufficient permissions.",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(mapToResponseDTO(user))
}

func (h *Handler) GetByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing email parameter",
		})
	}

	loggedInEmail := c.Locals("user_email")
	loggedInRole := c.Locals("user_role")
	if loggedInEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if loggedInRole != "admin" && loggedInEmail.(string) != email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden. Insufficient permissions.",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	user, err := h.service.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(mapToResponseDTO(user))
}

func (h *Handler) List(c *fiber.Ctx) error {
	loggedInRole := c.Locals("user_role")
	if loggedInRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden.",
		})
	}

	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

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

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	usersList, err := h.service.List(ctx, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users",
		})
	}

	response := make([]ResponseDTO, len(usersList))
	for i, u := range usersList {
		response[i] = mapToResponseDTO(&u)
	}

	return c.JSON(response)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing user ID parameter",
		})
	}

	loggedInID := c.Locals("user_id")
	loggedInRole := c.Locals("user_role")
	if loggedInID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if loggedInRole != "admin" && loggedInID.(string) != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden. Insufficient permissions.",
		})
	}

	dto := new(UpdateUserDTO)
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

	existingUser, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if dto.Email != nil {
		existingUser.Email = *dto.Email
	}
	if dto.FirstName != nil {
		existingUser.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		existingUser.LastName = *dto.LastName
	}
	if dto.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		existingUser.PasswordHash = string(hashed)
	}
	if dto.ProfileImageURL != nil {
		existingUser.ProfileImageURL = dto.ProfileImageURL
	}
	if dto.Provider != nil {
		existingUser.Provider = *dto.Provider
	}
	if dto.ProviderID != nil {
		existingUser.ProviderID = dto.ProviderID
	}
	if dto.Role != nil {
		if loggedInRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden.",
			})
		}
		existingUser.Role = *dto.Role
	}

	updatedUser, err := h.service.Update(ctx, existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(mapToResponseDTO(updatedUser))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing user ID parameter",
		})
	}

	loggedInID := c.Locals("user_id")
	loggedInRole := c.Locals("user_role")
	if loggedInID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if loggedInRole != "admin" && loggedInID.(string) != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden. Insufficient permissions.",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if _, err := h.service.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

func mapToResponseDTO(u *User) ResponseDTO {
	return ResponseDTO{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		ProfileImageURL: u.ProfileImageURL,
		Provider:        u.Provider,
		ProviderID:      u.ProviderID,
		Role:            u.Role,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
