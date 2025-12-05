package routes

import (
	"github.com/Lutfania/ekrp/app/repository"
	"github.com/Lutfania/ekrp/app/service"
	"github.com/Lutfania/ekrp/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	

	// Repositories
	userRepo := repository.NewUserRepository()
	achRepo := repository.NewAchievementRepository()
	studentRepo := repository.NewStudentRepository()
	lecturerRepo := repository.NewLecturerRepository()

	// Services
	authService := service.NewAuthService(userRepo)
	achService := service.NewAchievementService(achRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)
	lecturerService := service.NewLecturerService(lecturerRepo)

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

	
ach := app.Group("/api/v1/achievements", middleware.JWTAuth)

ach.Get("/", achService.List)
ach.Get("/:id", achService.GetByID)
ach.Post("/", achService.Create)
ach.Put("/:id", achService.Update)
ach.Delete("/:id", achService.Delete)

ach.Post("/:id/submit", achService.Submit)
ach.Post("/:id/verify", achService.Verify)
ach.Post("/:id/reject", achService.Reject)
ach.Get("/:id/history", achService.History)

	// STUDENTS
	students := app.Group("/api/v1/students", middleware.JWTAuth)
	students.Get("/", studentService.FindAll)
	students.Get("/:id", studentService.FindById)
	students.Post("/", studentService.Create)
	students.Put("/:id/advisor", studentService.UpdateAdvisor)
	students.Get("/:id/achievements", studentService.FindAchievements)

	// LECTURERS
	lecturers := app.Group("/api/v1/lecturers", middleware.JWTAuth)
	lecturers.Get("/", lecturerService.FindAll)
	lecturers.Get("/:id", lecturerService.FindById)
	lecturers.Post("/", lecturerService.Create)
	lecturers.Get("/:id/advisees", lecturerService.FindAdvisees)
}
