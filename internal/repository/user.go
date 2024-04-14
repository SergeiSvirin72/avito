package repository

import (
	"database/sql"
	"fmt"

	"avito/internal/model"
)

type User interface {
	GetByID(id int64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type user struct {
	db *sql.DB
}

func NewUser(db *sql.DB) User {
	return &user{
		db: db,
	}
}

func (r *user) GetByID(id int64) (*model.User, error) {
	var u model.User

	row := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role); err != nil {
		return nil, fmt.Errorf("userRepository.GetByID %d: %w", id, err)
	}

	return &u, nil
}

func (r *user) GetByEmail(email string) (*model.User, error) {
	var u model.User

	row := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role); err != nil {
		return nil, fmt.Errorf("userRepository.GetByEmail %s: %w", email, err)
	}

	return &u, nil
}
