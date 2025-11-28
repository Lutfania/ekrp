package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Ambil claims dari JWTAuth
		claims := c.Locals("claims")
		if claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		mapClaims := claims.(map[string]interface{})

		// Ambil permissions
		rawPerms := mapClaims["permissions"]
		if rawPerms == nil {
			return c.Status(403).JSON(fiber.Map{"error": "No permissions in token"})
		}

		perms := rawPerms.([]interface{})

		has := false
		for _, p := range perms {
			if p.(string) == permission {
				has = true
				break
			}
		}

		if !has {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: insufficient permissions"})
		}

		return c.Next()
	}
}
