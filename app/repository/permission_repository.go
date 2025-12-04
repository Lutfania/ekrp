package repository

import (
	"context"
	"github.com/Lutfania/ekrp/config"
)

type PermissionRepository struct{}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{}
}

func (r *PermissionRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := config.DB.Query(context.Background(), query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []string{}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		perms = append(perms, name)
	}

	return perms, nil
}
