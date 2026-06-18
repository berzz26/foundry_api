package auth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var RolePrivileges = map[string][]string{}

func LoadRolePrivileges(ctx context.Context, db *pgxpool.Pool) error {
	rows, err := db.Query(ctx, "SELECT role::text, privilege FROM role_privileges")
	if err != nil {
		return err
	}
	defer rows.Close()

	newPrivs := map[string][]string{}
	for rows.Next() {
		var role, privilege string
		if err := rows.Scan(&role, &privilege); err != nil {
			return err
		}
		newPrivs[role] = append(newPrivs[role], privilege)
	}
	RolePrivileges = newPrivs
	return nil
}

func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if os.Getenv("AUTH_BYPASS") == "true" {
			c.Locals("user_id", "00000000-0000-0000-0000-000000000000")
			c.Locals("user_email", "bypass-dev-user@example.com")
			c.Locals("user_role", "admin")
			c.Locals("user_privileges", RolePrivileges["admin"])
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
			c.Locals("user_privileges", RolePrivileges[claims.Role])
		}

		return c.Next()
	}
}

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Next()
	}
}

func RequirePrivilege(requiredPrivilege string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		privs := c.Locals("user_privileges")
		if privs == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		privList, ok := privs.([]string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error. Invalid privileges context.",
			})
		}

		for _, p := range privList {
			if p == requiredPrivilege {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}

func HasPrivilege(c *fiber.Ctx, privilege string) bool {
	privs := c.Locals("user_privileges")
	if privs == nil {
		return false
	}
	privList, ok := privs.([]string)
	if !ok {
		return false
	}
	for _, p := range privList {
		if p == privilege {
			return true
		}
	}
	return false
}
