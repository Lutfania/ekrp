package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/config"
)

type AchievementRepository struct{}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

func (r *AchievementRepository) Create(ar *models.AchievementReference) error {
	query := `INSERT INTO achievement_references
	(id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := config.DB.Exec(context.Background(), query,
		ar.StudentID, ar.MongoAchievementID, ar.Status,
		nil, nil, nil, ar.RejectionNote,
		ar.CreatedAt, ar.UpdatedAt,
	)
	return err
}

func (r *AchievementRepository) FindByID(id string) (*models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	 FROM achievement_references WHERE id = $1 LIMIT 1`
	row := config.DB.QueryRow(context.Background(), query, id)
	ar := &models.AchievementReference{}
	var submittedAt, verifiedAt sql.NullTime
	var verifiedBy, rejectionNote sql.NullString
	var updatedAt sql.NullTime
	err := row.Scan(&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
		&submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ar.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if submittedAt.Valid {
		ar.SubmittedAt = &submittedAt.Time
	}
	if verifiedAt.Valid {
		ar.VerifiedAt = &verifiedAt.Time
	}
	if verifiedBy.Valid {
		v := verifiedBy.String
		ar.VerifiedBy = &v
	}
	if rejectionNote.Valid {
		n := rejectionNote.String
		ar.RejectionNote = &n
	}
	if updatedAt.Valid {
		ar.UpdatedAt = &updatedAt.Time
	}
	return ar, nil
}

func (r *AchievementRepository) ListAll() ([]models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at FROM achievement_references ORDER BY created_at DESC`
	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.AchievementReference
	for rows.Next() {
		var ar models.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy, rejectionNote sql.NullString
		var updatedAt sql.NullTime
		if err := rows.Scan(&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ar.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if submittedAt.Valid {
			ar.SubmittedAt = &submittedAt.Time
		}
		if verifiedAt.Valid {
			ar.VerifiedAt = &verifiedAt.Time
		}
		if verifiedBy.Valid {
			v := verifiedBy.String
			ar.VerifiedBy = &v
		}
		if rejectionNote.Valid {
			n := rejectionNote.String
			ar.RejectionNote = &n
		}
		if updatedAt.Valid {
			ar.UpdatedAt = &updatedAt.Time
		}
		res = append(res, ar)
	}
	return res, nil
}

func (r *AchievementRepository) ListByStudent(studentID string) ([]models.AchievementReference, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at FROM achievement_references WHERE student_id=$1 ORDER BY created_at DESC`
	rows, err := config.DB.Query(context.Background(), query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.AchievementReference
	for rows.Next() {
		var ar models.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy, rejectionNote sql.NullString
		var updatedAt sql.NullTime
		if err := rows.Scan(&ar.ID, &ar.StudentID, &ar.MongoAchievementID, &ar.Status,
			&submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ar.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if submittedAt.Valid {
			ar.SubmittedAt = &submittedAt.Time
		}
		if verifiedAt.Valid {
			ar.VerifiedAt = &verifiedAt.Time
		}
		if verifiedBy.Valid {
			v := verifiedBy.String
			ar.VerifiedBy = &v
		}
		if rejectionNote.Valid {
			n := rejectionNote.String
			ar.RejectionNote = &n
		}
		if updatedAt.Valid {
			ar.UpdatedAt = &updatedAt.Time
		}
		res = append(res, ar)
	}
	return res, nil
}

func (r *AchievementRepository) UpdateStatus(id, status string, submittedAt, verifiedAt *time.Time, verifiedBy *string, rejectionNote *string) error {
	query := `UPDATE achievement_references SET status=$1, submitted_at=$2, verified_at=$3, verified_by=$4, rejection_note=$5, updated_at=$6 WHERE id=$7`
	_, err := config.DB.Exec(context.Background(), query, status, submittedAt, verifiedAt, verifiedBy, rejectionNote, time.Now(), id)
	return err
}

func (r *AchievementRepository) UpdateMongoID(id, mongoID string) error {
	query := `UPDATE achievement_references SET mongo_achievement_id=$1, updated_at=$2 WHERE id=$3`
	_, err := config.DB.Exec(context.Background(), query, mongoID, time.Now(), id)
	return err
}

func (r *AchievementRepository) Delete(id string) error {
	_, err := config.DB.Exec(context.Background(), `DELETE FROM achievement_references WHERE id=$1`, id)
	return err
}
