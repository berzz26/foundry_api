package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if os.Getenv("AUTH_BYPASS") == "true" {
			c.Locals("user_id", "00000000-0000-0000-0000-000000000000")
			c.Locals("user_email", "bypass-dev-user@example.com")
			c.Locals("user_role", "admin")
			return c.Next()
		}

		var tokenStr string

		// 1. Check Authorization header
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Check cookie
		if tokenStr == "" {
			tokenStr = c.Cookies("__Secure-token")
		}
		if tokenStr == "" {
			tokenStr = c.Cookies("token")
		}

		if tokenStr == "" {
			return c.Next()
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// Reject none signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err == nil && token.Valid {
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("user_role", claims.Role)
		}

		return c.Next()
	}
}

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized. Please sign in to access this resource.",
			})
		}
		return c.Next()
	}
}

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("user_role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized. Role validation failed.",
			})
		}

		roleStr, ok := role.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error. Invalid role type.",
			})
		}

		for _, r := range allowedRoles {
			if strings.EqualFold(r, roleStr) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden. Insufficient permissions.",
		})
	}
}
