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
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive)
	return err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, username, email, full_name, role_id, is_active FROM users`)
	if err != nil {
		return nil, err
	}

	var result []models.User

	for rows.Next() {
		u := models.User{}
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive)
		result = append(result, u)
	}

	return result, nil
}

func (r *UserRepository) FindById(id string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(),
		`SELECT id, username, email, full_name, role_id, is_active 
		 FROM users WHERE id = $1`, id)

	u := models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(),
		`SELECT id, username, email, password_hash, full_name, role_id, is_active
		 FROM users WHERE email=$1`, email)

	u := models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetRolePermissions(roleID string) ([]string, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT p.name 
		 FROM role_permissions rp
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE rp.role_id = $1`, roleID)

	if err != nil {
		return nil, err
	}

	var permissions []string
	for rows.Next() {
		var permName string
		if err := rows.Scan(&permName); err != nil {
			return nil, err
		}
		permissions = append(permissions, permName)
	}

	return permissions, nil
}



func (r *UserRepository) UpdateUser(id string, req *models.UpdateUserRequest) error {
	_, err := config.DB.Exec(context.Background(),
		`UPDATE users SET username=$1, email=$2, full_name=$3 WHERE id=$4`,
		req.Username, req.Email, req.FullName, id)
	return err
}

func (r *UserRepository) DeleteUser(id string) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepository) UpdateUserRole(id, roleID string) error {
	_, err := config.DB.Exec(context.Background(),
		`UPDATE users SET role_id=$1 WHERE id=$2`, roleID, id)
	return err
}
