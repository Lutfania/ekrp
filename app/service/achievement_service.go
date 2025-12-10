package service

import (
	"context"
	"io"
	"time"

	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/app/repository"
	"github.com/Lutfania/ekrp/config"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// AchievementService menangani logic yg gabungkan Postgres (reference) dan Mongo (dokumen prestasi)
type AchievementService struct {
	PGRepo    *repository.AchievementRepository
	MongoRepo *repository.MongoAchievementRepository
}

func NewAchievementService(pg *repository.AchievementRepository, mongo *repository.MongoAchievementRepository) *AchievementService {
	return &AchievementService{PGRepo: pg, MongoRepo: mongo}
}

// List -> GET /api/v1/achievements?student_id=...
func (s *AchievementService) List(c *fiber.Ctx) error {
	roleRaw := c.Locals("role_id")
	studentIDQuery := c.Query("student_id")

	// simplify: if admin (role name/id "Admin") => list all or filter by student
	if roleStr, ok := roleRaw.(string); ok && roleStr == "Admin" {
		if studentIDQuery != "" {
			list, err := s.PGRepo.ListByStudent(studentIDQuery)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			return s.buildAchievementResponses(list, true)
		}
		list, err := s.PGRepo.ListAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return s.buildAchievementResponses(list, true)
	}

	// non-admin (mahasiswa) — require student_id param (or adapt mapping user->student)
	if studentIDQuery == "" {
		return c.Status(400).JSON(fiber.Map{"error": "student_id required"})
	}
	list, err := s.PGRepo.ListByStudent(studentIDQuery)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return s.buildAchievementResponses(list, false)
}

// helper: build responses merging mongo doc
func (s *AchievementService) buildAchievementResponses(list []models.AchievementReference, includeAll bool) error {
	// not used; placeholder to satisfy signature if used elsewhere
	_ = includeAll
	return nil
}

// buildAchievementResponses returns JSON result — implementation returns in caller instead of here
func (s *AchievementService) buildAchievementResponsesWithData(list []models.AchievementReference) ([]models.AchievementResponse, error) {
	var out []models.AchievementResponse
	for _, ar := range list {
		resp := models.AchievementResponse{
			ID:                 ar.ID,
			StudentID:          ar.StudentID,
			MongoAchievementID: ar.MongoAchievementID,
			Status:             ar.Status,
			SubmittedAt:        ar.SubmittedAt,
			VerifiedAt:         ar.VerifiedAt,
			VerifiedBy:         ar.VerifiedBy,
			RejectionNote:      ar.RejectionNote,
			CreatedAt:          ar.CreatedAt,
			UpdatedAt:          ar.UpdatedAt,
		}

		// try fetch mongo doc if exists
		if ar.MongoAchievementID != "" {
			doc, err := s.MongoRepo.FindByIDHex(ar.MongoAchievementID)
			if err == nil && doc != nil {
				// adapt doc into map[string]interface{} for response
				m := map[string]interface{}{
					"id":          doc.ID,
					"title":       doc.Title,
					"description": doc.Description,
					"files":       doc.Files,
					"extra":       doc.Extra,
					"created_at":  doc.CreatedAt,
					"updated_at":  doc.UpdatedAt,
				}
				resp.Doc = m
			}
		}
		out = append(out, resp)
	}
	return out, nil
}

// GetByID -> GET /api/v1/achievements/:id
func (s *AchievementService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ar, err := s.PGRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	resp := models.AchievementResponse{
		ID:                 ar.ID,
		StudentID:          ar.StudentID,
		MongoAchievementID: ar.MongoAchievementID,
		Status:             ar.Status,
		SubmittedAt:        ar.SubmittedAt,
		VerifiedAt:         ar.VerifiedAt,
		VerifiedBy:         ar.VerifiedBy,
		RejectionNote:      ar.RejectionNote,
		CreatedAt:          ar.CreatedAt,
		UpdatedAt:          ar.UpdatedAt,
	}
	if ar.MongoAchievementID != "" {
		doc, err := s.MongoRepo.FindByIDHex(ar.MongoAchievementID)
		if err == nil && doc != nil {
			resp.Doc = map[string]interface{}{
				"id":          doc.ID,
				"title":       doc.Title,
				"description": doc.Description,
				"files":       doc.Files,
				"extra":       doc.Extra,
				"created_at":  doc.CreatedAt,
				"updated_at":  doc.UpdatedAt,
			}
		}
	}
	return c.JSON(resp)
}

// Create -> POST /api/v1/achievements
// expects models.CreateAchievementRequest in models (Doc map[string]interface{})
func (s *AchievementService) Create(c *fiber.Ctx) error {
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "student_id required"})
	}

	// build mongo document from req.Doc; we'll store under Extra field
	mongoDoc := &models.MongoAchievement{
		StudentID:   req.StudentID,
		Extra:       req.Doc,
		CreatedAt:   time.Now(),
	}
	hexID, err := s.MongoRepo.Insert(mongoDoc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	now := time.Now()
	ar := &models.AchievementReference{
		StudentID:          req.StudentID,
		MongoAchievementID: hexID,
		Status:             "draft",
		CreatedAt:          now,
	}
	if err := s.PGRepo.Create(ar); err != nil {
		// attempt cleanup in mongo (best effort)
		_ = s.MongoRepo.DeleteByHex(hexID)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "created", "mongo_id": hexID})
}

// Update -> PUT /api/v1/achievements/:id
func (s *AchievementService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	ar, err := s.PGRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	// update mongo doc if provided
	if req.MongoAchievementID != nil && *req.MongoAchievementID != "" {
		// update PG record's mongo id
		if err := s.PGRepo.UpdateMongoID(id, *req.MongoAchievementID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		ar.MongoAchievementID = *req.MongoAchievementID
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

// Delete -> DELETE /api/v1/achievements/:id
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	ar, err := s.PGRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	// delete mongo doc if exist
	if ar.MongoAchievementID != "" {
		_ = s.MongoRepo.DeleteByHex(ar.MongoAchievementID)
	}
	// delete pg reference
	if err := s.PGRepo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// Submit -> POST /api/v1/achievements/:id/submit
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	now := time.Now()
	if err := s.PGRepo.UpdateStatus(id, "submitted", &now, nil, nil, nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// optional history: insert to history table if exists (best effort)
	_ = insertHistoryIfTableExists(id, "draft", "submitted", c.Locals("user_id"))
	return c.JSON(fiber.Map{"message": "submitted"})
}

// Verify -> POST /api/v1/achievements/:id/verify
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	now := time.Now()
	verifier, _ := c.Locals("user_id").(string)
	if verifier == "" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	if err := s.PGRepo.UpdateStatus(id, "verified", nil, &now, &verifier, nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	_ = insertHistoryIfTableExists(id, "submitted", "verified", verifier)
	return c.JSON(fiber.Map{"message": "verified"})
}

// Reject -> POST /api/v1/achievements/:id/reject
func (s *AchievementService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	var body models.RejectRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	verifier, _ := c.Locals("user_id").(string)
	if verifier == "" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	if err := s.PGRepo.UpdateStatus(id, "rejected", nil, nil, &verifier, &body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	_ = insertHistoryIfTableExists(id, "submitted", "rejected", verifier)
	return c.JSON(fiber.Map{"message": "rejected"})
}

// History -> GET /api/v1/achievements/:id/history
func (s *AchievementService) History(c *fiber.Ctx) error {
	id := c.Params("id")
	// try reading dedicated history table; fallback to returning single reference
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, old_status, new_status, changed_by, note, changed_at
		 FROM achievement_reference_history WHERE achievement_ref_id=$1 ORDER BY changed_at DESC`, id)
	if err != nil {
		// fallback: return current record
		ar, err2 := s.PGRepo.FindByID(id)
		if err2 != nil {
			return c.Status(500).JSON(fiber.Map{"error": err2.Error()})
		}
		return c.JSON([]models.AchievementReference{*ar})
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var hid, oldS, newS string
		var changedBy, note *string
		var changedAt time.Time
		rows.Scan(&hid, &oldS, &newS, &changedBy, &note, &changedAt)
		history = append(history, map[string]interface{}{
			"id":         hid, "old_status": oldS, "new_status": newS,
			"changed_by": changedBy, "note": note, "changed_at": changedAt,
		})
	}
	return c.JSON(history)
}

// UploadAttachment -> POST /api/v1/achievements/:id/attachments
// multipart/form-data; field "file" (single)
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	ar, err := s.PGRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}
	mongoHex := ar.MongoAchievementID
	if mongoHex == "" {
		return c.Status(400).JSON(fiber.Map{"error": "no mongo document linked"})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}

	// read file content if needed (here we won't store to disk; just metadata). In production, upload to storage (S3) and save URL.
	f, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot open file"})
	}
	defer f.Close()
	// read small preview or skip reading large file; we'll discard content
	_, _ = io.Copy(io.Discard, f)

	fileMeta := map[string]interface{}{
		"file_name":   fileHeader.Filename,
		"file_size":   fileHeader.Size,
		"content_type": fileHeader.Header.Get("Content-Type"),
		"uploaded_at": time.Now(),
		// "file_url": "https://... if you upload to storage"
	}

	// push into mongo "files" array
	if err := s.MongoRepo.UpdateByHex(mongoHex, bson.M{"$push": bson.M{"files": fileMeta}}); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "attachment uploaded"})
}

/*** small helper ***/
func insertHistoryIfTableExists(achievementRefID, oldStatus, newStatus string, changedBy interface{}) error {
	// best-effort insert; if table doesn't exist it will error and we ignore
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO achievement_reference_history (id, achievement_ref_id, old_status, new_status, changed_by, note, changed_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, now())`, achievementRefID, oldStatus, newStatus, changedBy, nil)
	if err != nil {
		// ignore error (table may not exist)
		return err
	}
	return nil
}
