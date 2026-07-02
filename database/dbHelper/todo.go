package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/models"
	"database/sql"
	"errors"
	"strings"
)

func TodoExists(name, userID string) (bool, error) {

	query := `SELECT COUNT(id) > 0 FROM todos 
              WHERE TRIM(LOWER(name)) = TRIM(LOWER($1)) 
              AND user_id = $2
              AND archived_at IS NULL`

	var exists bool

	err := database.DB.Get(&exists, query, name, userID)

	return exists, err
}
func GetAllTodos(userID string) ([]models.Todo, error) {
	todos := make([]models.Todo, 0)
	query := `SELECT id,name,description,is_completed,created_at, archived_at 
            FROM  todos 
            WHERE user_id=$1 
            AND archived_at IS NULL
            ORDER BY id ASC`
	err := database.DB.Select(&todos, query, userID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func CreateTodo(req models.TodoRequest) (models.Todo, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	args := []interface{}{
		req.Name,
		req.Description,
		req.UserID,
	}
	var todo models.Todo
	query := `
		INSERT INTO todos (name, description,user_id)
		VALUES ($1, $2,$3)
		RETURNING id, name, description, is_completed,user_id, created_at,  archived_at`

	err := database.DB.QueryRowx(query, args...).StructScan(&todo)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Todo{}, nil
	}
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func GetTodoById(id, userID string) (models.Todo, error) {
	var todo models.Todo
	query := `Select *
            from todos 
            where id=$1
            AND user_id = $2
            AND archived_at IS NULL`

	err := database.DB.QueryRowx(query, id, userID).StructScan(&todo)

	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}
func UpdateTodo(id, userID string, req models.TodoRequest) (models.Todo, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	var todo models.Todo
	args := []interface{}{
		id,
		userID,
		req.Name,
		req.Description,
	}

	query := `
		UPDATE todos
		SET name=$3,
		    description=$4
		WHERE id=$1
		AND user_id=$2
		AND archived_at IS NULL 
		RETURNING id, name, description, is_completed,user_id,created_at,archived_at
	`
	err := database.DB.QueryRowx(query, args...).StructScan(&todo)

	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}
func DeleteTodoById(id, userID string) (bool, error) {
	query := `UPDATE todos
            SET archived_at = NOW()
            WHERE id = $1
            AND user_id = $2
            AND archived_at IS NULL`
	result, err := database.DB.Exec(query, id, userID)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false, nil // ID nahi mili
	}
	return true, nil
}

func MarkTodoAsCompleted(id, userID string) (models.Todo, error) {
	var todo models.Todo

	query := `
		UPDATE todos
		SET is_completed = true
		WHERE id = $1
		AND user_id = $2
        AND archived_at IS NULL
		RETURNING id, name, description, is_completed,user_id,created_at,archived_at
	`

	err := database.DB.QueryRowx(query, id, userID).StructScan(&todo)

	if err != nil {
		return models.Todo{}, err

	}

	return todo, nil
}
func DeleteAllTodos(userID string) error {
	SQL := `UPDATE todos
			  SET archived_at = NOW()
			  WHERE user_id = $1
			    AND archived_at IS NULL`
	_, err := database.DB.Exec(SQL, userID)
	return err
}
