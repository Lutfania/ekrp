package service

import (
	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	Repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) FindAll(c *fiber.Ctx) error {
	list, err := s.Repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (s *StudentService) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	st, err := s.Repo.FindById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "student not found"})
	}
	return c.JSON(st)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	var req models.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.UserID == "" || req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id and student_id are required"})
	}
	if err := s.Repo.Create(&req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "student created"})
}

func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.AdvisorID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "advisor_id required"})
	}
	if err := s.Repo.UpdateAdvisor(id, req.AdvisorID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "advisor updated"})
}

func (s *StudentService) FindAchievements(c *fiber.Ctx) error {
	id := c.Params("id")
	list, err := s.Repo.FindAchievements(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}
