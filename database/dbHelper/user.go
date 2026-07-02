package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Querier interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func CreateUser(tx *sqlx.Tx, name, email, hashedPassword string) (string, error) {
	var userID string
	query := `INSERT INTO users (name, email, password)
              VALUES ($1, $2, $3)
              RETURNING id`
	err := tx.QueryRow(query, strings.TrimSpace(name), strings.TrimSpace(strings.ToLower(email)), hashedPassword).Scan(&userID)
	return userID, err
}

func IsUSerExist(email string) (bool, error) {
	query := `SELECT count(id)>0
            FROM users where email=TRIM($1)
            AND archived_at is null `
	var check bool
	err := database.DB.Get(&check, query, email)
	return check, err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT *
              FROM users 
              WHERE email=TRIM(LOWER($1)) 
              AND archived_at IS NULL`
	err := database.DB.Get(&user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUserSession(db Querier, userID string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session (user_id)
              VALUES ($1)
              RETURNING id`
	err := db.QueryRow(query, userID).Scan(&sessionID)
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

func DeleteUser(db sqlx.Execer, userID string) error {
	query := `UPDATE users
              SET archived_at = NOW()
              WHERE id = $1
                AND archived_at IS NULL`
	_, err := db.Exec(query, userID)
	return err
}

func GetUserBySession(sessionID string) (*models.User, error) {
	var user models.User
	SQL := `SELECT u.id, u.name, u.email
			  FROM users u
			  JOIN user_session s ON s.user_id = u.id
			  WHERE s.id = $1
			    AND s.archived_at IS NULL
			    AND u.archived_at IS NULL`
	err := database.DB.Get(&user, SQL, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
