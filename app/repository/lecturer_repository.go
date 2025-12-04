package repository

import (
	"context"
	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/config"
)

type LecturerRepository struct{}

func NewLecturerRepository() *LecturerRepository {
	return &LecturerRepository{}
}

// FindAll lecturers
func (r *LecturerRepository) FindAll() ([]models.Lecturer, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, user_id, lecturer_id, department, created_at FROM lecturers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Lecturer
	for rows.Next() {
		l := models.Lecturer{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, l)
	}
	return res, nil
}

// FindById lecturer
func (r *LecturerRepository) FindById(id string) (*models.Lecturer, error) {
	row := config.DB.QueryRow(context.Background(),
		`SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE id = $1 LIMIT 1`, id)
	l := &models.Lecturer{}
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		return nil, err
	}
	return l, nil
}

// Create a lecturer
func (r *LecturerRepository) Create(l *models.Lecturer) error {
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO lecturers (user_id, lecturer_id, department) VALUES ($1, $2, $3)`,
		l.UserID, l.LecturerID, l.Department)
	return err
}

// Find advisees (students) by lecturer id (returns student rows)
func (r *LecturerRepository) FindAdvisees(lecturerID string) ([]map[string]interface{}, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		 FROM students WHERE advisor_id = $1`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []map[string]interface{}
	for rows.Next() {
		var id, userID, studentID, programStudy, academicYear, advisorID string
		var createdAt any
		if err := rows.Scan(&id, &userID, &studentID, &programStudy, &academicYear, &advisorID, &createdAt); err != nil {
			return nil, err
		}
		rec := map[string]interface{}{
			"id":            id,
			"user_id":       userID,
			"student_id":    studentID,
			"program_study": programStudy,
			"academic_year": academicYear,
			"advisor_id":    advisorID,
			"created_at":    createdAt,
		}
		res = append(res, rec)
	}
	return res, nil
}
