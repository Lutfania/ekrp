package service

import (
	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/app/repository"
	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	Repo *repository.LecturerRepository
}

func NewLecturerService(repo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{Repo: repo}
}

// GET /api/v1/lecturers
func (s *LecturerService) FindAll(c *fiber.Ctx) error {
	list, err := s.Repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// GET /api/v1/lecturers/:id
func (s *LecturerService) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	lect, err := s.Repo.FindById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Lecturer not found"})
	}
	return c.JSON(lect)
}

// POST /api/v1/lecturers
func (s *LecturerService) Create(c *fiber.Ctx) error {
	var req struct {
		UserID     string `json:"user_id"`
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.UserID == "" || req.LecturerID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id and lecturer_id required"})
	}

	l := &models.Lecturer{
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}
	if err := s.Repo.Create(l); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Lecturer created"})
}

// GET /api/v1/lecturers/:id/advisees
func (s *LecturerService) FindAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")
	rows, err := s.Repo.FindAdvisees(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rows)
}
