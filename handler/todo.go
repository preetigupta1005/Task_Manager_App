package handler

import (
	"My-todo-app/database/dbHelper"
	"My-todo-app/model"
	"My-todo-app/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := dbHelper.GetAllTodos()
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req model.TodoRequest

	err := utils.ParseBody(r, &req)
	if err != nil {
		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	exists, err := dbHelper.TodoExists(req.Name)
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
	idstr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idstr)
	todo, err := dbHelper.GetTodoById(id)
	if err != nil {
		log.Println("Error:", err)
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var req model.TodoRequest
	err := utils.ParseBody(r, &req)
	if err != nil {
		return
	}
	todo, err := dbHelper.UpdateTodo(id, req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}
	utils.RespondJSON(w, http.StatusOK, todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	found, err := dbHelper.DeleteTodoById(id)
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
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	todo, err := dbHelper.MarkTodoAsCompleted(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}
