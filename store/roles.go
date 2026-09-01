package store

import (
	"context"
	"database/sql"
)

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}
type RolesStorage struct {
	db *sql.DB
}

func (s *RolesStorage) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	query := `
				SELECT * FROM roles WHERE name=$1
		`

	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&role.Id,
		&role.Name,
		&role.Level,
		&role.Description,
	)
	if err != nil {
		return nil, ErrorFactoryDB(err)
	}
	return &role, nil
}
