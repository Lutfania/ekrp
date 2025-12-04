package middleware

import (
	"github.com/Lutfania/ekrp/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing authorization"})
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
	}

	claims, err := utils.ValidateToken(parts[1])
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	c.Locals("claims", claims)
	c.Locals("user_id", claims.UserID)
	c.Locals("role_id", claims.RoleID)

	return c.Next()
}
