package users

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user User) (*User, error) {
	row := r.db.QueryRow(`INSERT INTO users (email, username, password_hash) VALUES($1, $2, $3) RETURNING id;`, user.Email, user.Username, user.Password)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("repository create user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Login(email string) (*User, error) {
	user := &User{
		Email: email,
	}
	row := r.db.QueryRow(`SELECT id, username, password_hash FROM users WHERE email=$1;`, user.Email)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("repository login: %w", err)
	}
	return user, nil
}
