package repository

import (
	"context"
	"time"

	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/config"
)

type AchievementRepository struct{}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

func (r *AchievementRepository) Create(ar *models.AchievementReference) error {
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW(), $8)`,
		ar.StudentID, ar.MongoAchievementID, ar.Status, ar.SubmittedAt, ar.VerifiedAt, ar.VerifiedBy, ar.RejectionNote, ar.UpdatedAt,
	)
	return err
}

func (r *AchievementRepository) FindByID(id string) (*models.AchievementReference, error) {
	row := config.DB.QueryRow(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		 FROM achievement_references WHERE id=$1 LIMIT 1`, id)

	ar := &models.AchievementReference{}
	var verifiedBy, rejectionNote *string
	var submittedAt, verifiedAt *time.Time
	var updatedAt *time.Time

	err := row.Scan(
		&ar.ID,
		&ar.StudentID,
		&ar.MongoAchievementID,
		&ar.Status,
		&submittedAt,
		&verifiedAt,
		&verifiedBy,
		&rejectionNote,
		&ar.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	ar.SubmittedAt = submittedAt
	ar.VerifiedAt = verifiedAt
	ar.VerifiedBy = verifiedBy
	ar.RejectionNote = rejectionNote
	ar.UpdatedAt = updatedAtValue(updatedAt)
	return ar, nil
}

func (r *AchievementRepository) ListAll() ([]models.AchievementReference, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at FROM achievement_references ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.AchievementReference
	for rows.Next() {
		ar := models.AchievementReference{}
		var verifiedBy, rejectionNote *string
		var submittedAt, verifiedAt *time.Time
		var updatedAt *time.Time
		if err := rows.Scan(&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status, &submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ar.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		ar.SubmittedAt = submittedAt
		ar.VerifiedAt = verifiedAt
		ar.VerifiedBy = verifiedBy
		ar.RejectionNote = rejectionNote
		ar.UpdatedAt = updatedAtValue(updatedAt)
		res = append(res, ar)
	}
	return res, nil
}

func (r *AchievementRepository) ListByStudent(studentID string) ([]models.AchievementReference, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		 FROM achievement_references WHERE student_id=$1 ORDER BY created_at DESC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.AchievementReference
	for rows.Next() {
		ar := models.AchievementReference{}
		var verifiedBy, rejectionNote *string
		var submittedAt, verifiedAt *time.Time
		var updatedAt *time.Time
		if err := rows.Scan(&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status, &submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ar.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		ar.SubmittedAt = submittedAt
		ar.VerifiedAt = verifiedAt
		ar.VerifiedBy = verifiedBy
		ar.RejectionNote = rejectionNote
		ar.UpdatedAt = updatedAtValue(updatedAt)
		res = append(res, ar)
	}
	return res, nil
}

func (r *AchievementRepository) UpdateStatus(id, status string, submittedAt, verifiedAt *time.Time, verifiedBy *string, rejectionNote *string) error {
	_, err := config.DB.Exec(context.Background(),
		`UPDATE achievement_references SET status=$1, submitted_at=$2, verified_at=$3, verified_by=$4, rejection_note=$5, updated_at=NOW() WHERE id=$6`,
		status, submittedAt, verifiedAt, verifiedBy, rejectionNote, id)
	return err
}

func (r *AchievementRepository) Delete(id string) error {
	_, err := config.DB.Exec(context.Background(), `DELETE FROM achievement_references WHERE id=$1`, id)
	return err
}

func updatedAtValue(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	return t
}
