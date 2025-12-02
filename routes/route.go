package routes

import (
	"ekrp/app/repository"
	"ekrp/app/service"
	"ekrp/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {

	// Repository
	userRepo := repository.NewUserRepository()

	// Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// AUTH
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.RefreshToken)
	auth.Post("/logout", authService.Logout)
	auth.Get("/profile", middleware.JWTAuth, authService.Profile)

	// USERS
	users := app.Group("/api/v1/users", middleware.JWTAuth)

	users.Get("/", userService.FindAll)
	users.Get("/:id", userService.FindById)
	users.Post("/", userService.CreateUser)
	users.Put("/:id", userService.UpdateUser)
	users.Delete("/:id", userService.DeleteUser)
	users.Put("/:id/role", userService.UpdateUserRole)
}
