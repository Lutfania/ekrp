package routes

import (
    "ekrp/app/models"
    "ekrp/app/repository"
    "ekrp/app/service"

    "github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

    // repository
    userRepo := repository.NewUserRepository()

    // service pakai constructor (lebih clean)
    authService := service.NewAuthService(userRepo)

    app.Post("/login", func(c *fiber.Ctx) error {
        var req models.LoginRequest

        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
        }

        res, err := authService.Login(req)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": err.Error()})
        }

        return c.JSON(res)
    })
}
