package repository

import (
	"context"
	"database/sql"
	"ekrp/app/models"
	"ekrp/config"
)

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{}
}

func (r *StudentRepository) FindAll() ([]models.Student, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Student
	for rows.Next() {
		var s models.Student
		var advisor sql.NullString
		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisor, &s.CreatedAt); err != nil {
			return nil, err
		}
		if advisor.Valid {
			val := advisor.String
			s.AdvisorID = &val
		} else {
			s.AdvisorID = nil
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *StudentRepository) FindById(id string) (*models.Student, error) {
	row := config.DB.QueryRow(context.Background(),
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE id = $1`, id)

	var s models.Student
	var advisor sql.NullString
	if err := row.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisor, &s.CreatedAt); err != nil {
		return nil, err
	}
	if advisor.Valid {
		val := advisor.String
		s.AdvisorID = &val
	} else {
		s.AdvisorID = nil
	}
	return &s, nil
}

func (r *StudentRepository) Create(req *models.CreateStudentRequest) error {
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
		 VALUES ($1, $2, $3, $4, $5)`,
		req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID)
	return err
}

func (r *StudentRepository) UpdateAdvisor(id string, advisorID string) error {
	_, err := config.DB.Exec(context.Background(),
		`UPDATE students SET advisor_id = $1 WHERE id = $2`, advisorID, id)
	return err
}

func (r *StudentRepository) FindAchievements(studentID string) ([]models.AchievementReference, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		 FROM achievement_references WHERE student_id = $1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.AchievementReference
	for rows.Next() {
		var a models.AchievementReference
		var submitted sql.NullTime
		var verified sql.NullTime
		var verifiedBy sql.NullString
		var rejection sql.NullString
		var updated sql.NullTime

		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status,
			&submitted, &verified, &verifiedBy, &rejection, &a.CreatedAt, &updated,
		); err != nil {
			return nil, err
		}
		if submitted.Valid {
			tmp := submitted.Time
			a.SubmittedAt = &tmp
		}
		if verified.Valid {
			tmp := verified.Time
			a.VerifiedAt = &tmp
		}
		if verifiedBy.Valid {
			val := verifiedBy.String
			a.VerifiedBy = &val
		}
		if rejection.Valid {
			val := rejection.String
			a.RejectionNote = &val
		}
		if updated.Valid {
			tmp := updated.Time
			a.UpdatedAt = &tmp
		}
		out = append(out, a)
	}
	return out, nil
}
