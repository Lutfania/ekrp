package middleware

import (
	"ekrp/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func JWTAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"missing authorization"})
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"invalid authorization"})
	}
	claims, err := utils.ValidateToken(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"invalid token"})
	}
	// simpan user id di Locals agar handler/service bisa ambil
	c.Locals("user_id", claims.UserID)
	c.Locals("role_id", claims.RoleID)
	return c.Next()
}
