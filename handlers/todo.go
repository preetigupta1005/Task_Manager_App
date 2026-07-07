package handlers

import (
	"My-todo-app/database"
	"My-todo-app/database/dbHelper"
	"My-todo-app/middlewares"
	"My-todo-app/models"
	"My-todo-app/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	todos, err := dbHelper.GetAllTodos(userCtx.UserID)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.TodoRequest
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	req.UserID = userCtx.UserID
	if err := utils.ParseBody(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if err := validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}
	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "name cannot be empty")
		return
	}

	exists, err := dbHelper.TodoExists(req.Name, req.UserID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "todo already exists", http.StatusConflict)
		return
	}

	todo, err := dbHelper.CreateTodo(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, todo)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	todo, err := dbHelper.GetTodoById(id, userCtx.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.RespondError(w, http.StatusNotFound, nil, "todo not found")
		return
	}
	if err != nil {
		log.Println("Error:", err)
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.TodoRequest
	err := utils.ParseBody(r, &req)
	if err := validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}
	if err != nil {
		return
	}
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	todo, err := dbHelper.UpdateTodo(id, userCtx.UserID, req)
	if errors.Is(err, sql.ErrNoRows) {
		utils.RespondError(w, http.StatusNotFound, nil, "todo not found")
		return
	}
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	utils.RespondJSON(w, http.StatusOK, todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	found, err := dbHelper.DeleteTodoById(id, userCtx.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func MarkTodoAsCompleted(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	todo, err := dbHelper.MarkTodoAsCompleted(id, userCtx.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to complete todo")

		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	if err := dbHelper.DeleteAllTodos(database.DB, userCtx.UserID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete all todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}
