package handler

import (
	"My-todo-app/database/dbHelper"
	"My-todo-app/middleware"
	"My-todo-app/model"
	"My-todo-app/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	todos, err := dbHelper.GetAllTodos(userCtx.UserID)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req model.TodoRequest
	userCtx := middleware.UserContext(r)
	req.UserID = userCtx.UserID
	if err := utils.ParseBody(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
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

	userCtx := middleware.UserContext(r)
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

	var req model.TodoRequest
	err := utils.ParseBody(r, &req)
	if err != nil {
		return
	}
	userCtx := middleware.UserContext(r)
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

	userCtx := middleware.UserContext(r)
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
	userCtx := middleware.UserContext(r)
	todo, err := dbHelper.MarkTodoAsCompleted(id, userCtx.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to complete todo")

		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}
