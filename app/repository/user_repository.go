package repository

import (
	"context"
	"ekrp/app/models"
	"ekrp/config"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := config.DB.Exec(context.Background(), query,
		user.Username,     // $1
		user.Email,        // $2
		user.PasswordHash, // $3
		user.FullName,     // $4  ← sudah benar
		user.RoleID,       // $5  ← ini yang penting!
		user.IsActive,     // $6
	)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := config.DB.QueryRow(context.Background(), query, email)

	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
