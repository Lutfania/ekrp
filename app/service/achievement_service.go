package service

import (
	"time"

	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/app/repository"
	"github.com/Lutfania/ekrp/utils"
    "context"
    "github.com/Lutfania/ekrp/config"

	"github.com/gofiber/fiber/v2"
)

type AchievementService struct {
	Repo *repository.AchievementRepository
}

func NewAchievementService(repo *repository.AchievementRepository) *AchievementService {
	return &AchievementService{Repo: repo}
}

// GET /api/v1/achievements
// - Admin sees all, mahasiswa sees only their own
func (s *AchievementService) List(c *fiber.Ctx) error {
	role := c.Locals("role_id")
_ = c.Locals("user_id")

	// optional query: ?student_id=...
	studentID := c.Query("student_id")

	if roleStr, ok := role.(string); ok && roleStr == /* admin role id? or name */ "Admin" {
		// admin => all (or can filter by student)
		if studentID != "" {
			list, err := s.Repo.ListByStudent(studentID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(list)
		}
		list, err := s.Repo.ListAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(list)
	}

	// non-admin => if student_id provided and matches, return; else get by user -> student mapping not included here
	// For simplicity: if student_id supplied, return that; otherwise require student_id
	if studentID == "" {
		// try to read student mapping from claims or ask client to pass student_id
		return c.Status(400).JSON(fiber.Map{"error": "student_id required"})
	}
	list, err := s.Repo.ListByStudent(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// GET /api/v1/achievements/:id
func (s *AchievementService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ar, err := s.Repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(ar)
}

// POST /api/v1/achievements  -> create draft
func (s *AchievementService) Create(c *fiber.Ctx) error {
	var req struct {
		StudentID          string `json:"student_id"`
		MongoAchievementID string `json:"mongo_achievement_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "student_id required"})
	}

	now := time.Now()
	ar := &models.AchievementReference{
		StudentID:          req.StudentID,
		MongoAchievementID: req.MongoAchievementID,
		Status:             "draft",
		CreatedAt:          now,
	}
	if err := s.Repo.Create(ar); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "created"})
}

// PUT /api/v1/achievements/:id  -> update (only certain fields)
func (s *AchievementService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		MongoAchievementID *string `json:"mongo_achievement_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// load, modify, save fields (simple approach: update only mongo id)
	ar, err := s.Repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if req.MongoAchievementID != nil {
		ar.MongoAchievementID = *req.MongoAchievementID
	}
	// store updated fields: use UpdateStatus with same status, no times changed
	if err := s.Repo.UpdateStatus(id, ar.Status, ar.SubmittedAt, ar.VerifiedAt, ar.VerifiedBy, ar.RejectionNote); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// DELETE /api/v1/achievements/:id
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.Repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// POST /api/v1/achievements/:id/submit  -> mahasiswa submit for verification
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	now := time.Now()
	status := "submitted"
	if err := s.Repo.UpdateStatus(id, status, &now, nil, nil, nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// optional: insert to history table (if created)
	_ = utils.InsertAchievementHistory(id, "draft", "submitted", c.Locals("user_id"))
	return c.JSON(fiber.Map{"message": "submitted"})
}

// POST /api/v1/achievements/:id/verify  -> dosen wali verify
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	now := time.Now()
	roleID, _ := c.Locals("role_id").(string)
	if roleID == "" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	verifier := c.Locals("user_id").(string)
	status := "verified"
	if err := s.Repo.UpdateStatus(id, status, nil, &now, &verifier, nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	_ = utils.InsertAchievementHistory(id, "submitted", "verified", verifier)
	return c.JSON(fiber.Map{"message": "verified"})
}

// POST /api/v1/achievements/:id/reject  -> dosen wali reject
func (s *AchievementService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	verifier := c.Locals("user_id").(string)
	status := "rejected"
	if err := s.Repo.UpdateStatus(id, status, nil, nil, &verifier, &body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	_ = utils.InsertAchievementHistory(id, "submitted", "rejected", verifier)
	return c.JSON(fiber.Map{"message": "rejected"})
}

// GET /api/v1/achievements/:id/history
func (s *AchievementService) History(c *fiber.Ctx) error {
	id := c.Params("id")
	// If you created table achievement_reference_history, query it:
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, old_status, new_status, changed_by, note, changed_at
		 FROM achievement_reference_history WHERE achievement_ref_id=$1 ORDER BY changed_at DESC`, id)
	if err != nil {
		// if no table, fallback: return current record only
		ar, err2 := s.Repo.FindByID(id)
		if err2 != nil {
			return c.Status(500).JSON(fiber.Map{"error": err2.Error()})
		}
		return c.JSON([]models.AchievementReference{*ar})
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var hid, oldS, newS string
		var changedBy *string
		var note *string
		var changedAt time.Time
		rows.Scan(&hid, &oldS, &newS, &changedBy, &note, &changedAt)
		history = append(history, map[string]interface{}{
			"id":         hid, "old_status": oldS, "new_status": newS,
			"changed_by": changedBy, "note": note, "changed_at": changedAt,
		})
	}
	return c.JSON(history)
}
