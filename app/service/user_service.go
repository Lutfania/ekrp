package service

import (
	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/app/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// GET /users
func (s *UserService) FindAll(c *fiber.Ctx) error {
	users, err := s.Repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GET /users/:id
func (s *UserService) FindById(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.Repo.FindById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// POST /users
func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User created"})
}

// PUT /users/:id
func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// FIX: harus pointer *
	if err := s.Repo.UpdateUser(id, &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated"})
}

// DELETE /users/:id
func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.Repo.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}

// PUT /users/:id/role
func (s *UserService) UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	type RoleUpdate struct {
		RoleID string `json:"role_id"`
	}

	var req RoleUpdate

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid role request"})
	}

	if err := s.Repo.UpdateUserRole(id, req.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role updated"})
}
