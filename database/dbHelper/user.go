package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateUser(name, email, hashedPassword string) (string, error) {
	var userID string
	query := `INSERT INTO users (name, email, password)
              VALUES ($1, $2, $3)
              RETURNING id`
	err := database.DB.QueryRow(query, strings.TrimSpace(name), strings.TrimSpace(strings.ToLower(email)), hashedPassword).Scan(&userID)
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

func CreateUserSession(db sqlx.Execer, userID, sessionToken string) error {
	SQL := `INSERT INTO user_session (user_id, session_token) VALUES ($1, $2)`
	_, err := db.Exec(SQL, userID, sessionToken)
	return err
}

func DeleteUserSession(db sqlx.Execer, sessionToken string) error {
	query := `UPDATE user_session
              SET archived_at = NOW()
              WHERE session_token  = $1
                AND archived_at IS NULL`
	_, err := db.Exec(query, sessionToken)
	return err
}

func DeleteUser(tx *sqlx.Tx, userID string) error {
	query := `UPDATE users
              SET archived_at = NOW()
              WHERE id = $1
                AND archived_at IS NULL`
	_, err := tx.Exec(query, userID)
	return err
}

func GetUser(userID string) (models.User, error) {
	var user models.User
	query := `Select id,name,email,created_at
            from users 
            where id=$1
            AND archived_at is null`

	err := database.DB.Get(&user, query, userID)
	return user, err

}
