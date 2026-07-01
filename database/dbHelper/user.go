package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/model"
	"My-todo-app/utils"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateUser(req model.RegisterRequest) error {
	hashedPassword, err := utils.HashedPassword(req.Password)
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
	query := `SELECT *
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

func CreateUserSession(userID string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id) 
              VALUES ($1) RETURNING id`
	err := database.DB.Get(&sessionID, query, userID)
	return sessionID, err
}

func DeleteUserSession(db sqlx.Execer, sessionID string) error {
	query := `UPDATE user_session
              SET archived_at = NOW()
              WHERE id = $1
                AND archived_at IS NULL`
	_, err := db.Exec(query, sessionID)
	return err
}

func IsSessionValid(sessionID string) (bool, error) {
	var isValid bool
	query := `SELECT count(id) > 0 
              FROM user_session 
              WHERE id = $1 
                AND archived_at IS NULL`
	err := database.DB.Get(&isValid, query, sessionID)
	return isValid, err
}

func DeleteUser(db sqlx.Execer, userID string) error {
	query := `UPDATE users
              SET archived_at = NOW()
              WHERE id = $1
                AND archived_at IS NULL`
	_, err := db.Exec(query, userID)
	return err
}

func GetUser(userID string) (model.User, error) {
	var user model.User
	SQL := `SELECT id, name, email 
              FROM users 
              WHERE id = $1
                AND archived_at IS NULL`
	getErr := database.DB.Get(&user, SQL, userID)
	return user, getErr
}
