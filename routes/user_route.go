package routes

import (
    "ekrp/app/service"
    "github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
    userService := service.NewUserService()

    app.Post("/users", func(c *fiber.Ctx) error {
        var body struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
            FullName string `json:"full_name"`
            RoleID   string `json:"role_id"`
        }

        if err := c.BodyParser(&body); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }

        err := userService.CreateUser(
            body.Username,
            body.Email,
            body.Password,
            body.FullName,
            body.RoleID,
        )

        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        return c.JSON(fiber.Map{"message": "User created successfully"})
    })
}
