package service

import (
	"ekrp/app/models"
	"ekrp/app/repository"
	"ekrp/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: repo,
	}
}

// ------------------------------------
// POST /auth/login
// ------------------------------------
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid email or password"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid email or password"})
	}

	// Ambil permissions dari role
	permissions, _ := s.UserRepo.GetRolePermissions(user.RoleID)

	token, err := utils.GenerateTokenWithPermissions(
		user.ID,
		user.RoleID,
		permissions,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(models.LoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		RoleID:      user.RoleID,
		Token:       token,
		Permissions: permissions,
	})
}

// ------------------------------------
// GET /auth/profile
// ------------------------------------
func (s *AuthService) Profile(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.JSON(user)
}

// ------------------------------------
// POST /auth/logout
// ------------------------------------
func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logged out"})
}

// ------------------------------------
// POST /auth/refresh
// ------------------------------------
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}
