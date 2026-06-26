package dbHelper

import (
	"My-todo-app/database"
	"My-todo-app/model"
	"database/sql"
	"errors"
	"strings"
)

func TodoExists(name string) (bool, error) {

	query := " SELECT COUNT(id) > 0 FROM todos WHERE name = TRIM($1)"

	var exists bool

	err := database.DB.Get(&exists, query, name)

	return exists, err
}
func GetAllTodos() ([]model.Todo, error) {
	todos := make([]model.Todo, 0)
	err := database.DB.Select(&todos, "SELECT id,name,description,is_completed FROM  todos ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func CreateTodo(req model.TodoRequest) (model.Todo, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	args := []interface{}{
		req.Name,
		req.Description,
	}
	var todo model.Todo
	query := `
		INSERT INTO todos (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, is_completed`

	err := database.DB.QueryRowx(query, args...).StructScan(&todo)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, nil
	}
	if err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}

func GetTodoById(id int) (model.Todo, error) {
	var todo model.Todo
	query := `Select id,name,description,is_completed 
            from todos 
            where id=$1`

	err := database.DB.QueryRowx(query, id).StructScan(&todo)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, nil
	}
	if err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}
func UpdateTodo(id int, req model.TodoRequest) (model.Todo, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	var todo model.Todo
	args := []interface{}{
		id,
		req.Name,
		req.Description,
	}

	query := `
		UPDATE todos
		SET name=$2,
		    description=$3
		WHERE id=$1
		RETURNING id, name, description, is_completed
	`
	err := database.DB.QueryRowx(query, args...).StructScan(&todo)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, nil
	}
	if err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}
func DeleteTodoById(id int) (bool, error) {
	result, err := database.DB.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false, nil // ID nahi mili
	}
	return true, nil
}

func MarkTodoAsCompleted(id int) (model.Todo, error) {
	var todo model.Todo

	query := `
		UPDATE todos
		SET is_completed = true
		WHERE id = $1
		RETURNING id, name, description, is_completed
	`

	err := database.DB.QueryRowx(query, id).StructScan(&todo)

	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}
