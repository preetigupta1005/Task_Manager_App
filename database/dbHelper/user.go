package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/model"
	"database/sql"
	"errors"
	"strings"
)

func CreateUser(req model.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	args := []interface{}{
		strings.TrimSpace(req.Name),
		strings.TrimSpace(strings.ToLower(req.Email)),
		hashedPassword,
	}
	query := `INSERT INTO users (name,email,password)
            VALUES ($1,$2,$3)  `
	_, insertErr := database.DB.Exec(query, args...)
	return insertErr
}

func IsUSerExist(email string) (bool, error) {
	query := `SELECT count(id)>0
            FROM users where email=TRIM($1)
            AND archived_at is null `
	var check bool
	err := database.DB.Get(&check, query, email)
	return check, err
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password 
              FROM users 
              WHERE email=TRIM(LOWER($1)) 
              AND archived_at IS NULL`
	err := database.DB.Get(&user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, nil
	}
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
