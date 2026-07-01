package handler

import (
	"My-todo-app/database"
	"My-todo-app/database/dbHelper"
	"My-todo-app/middleware"
	"My-todo-app/model"
	"My-todo-app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

var v = validator.New()

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userReq model.RegisterRequest
	err := utils.ParseBody(r, &userReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	if err := v.Struct(userReq); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	exists, err := dbHelper.IsUSerExist(userReq.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, nil, "user already exists")
		return
	}
	saveErr := dbHelper.CreateUser(userReq)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"user registered successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	if parseErr := utils.ParseBody(r, &req); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if err := v.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	user, err := dbHelper.GetUserByEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get user")
		return
	}

	if user.ID == "" {
		utils.RespondError(w, http.StatusNotFound, nil, "user not found")
		return
	}

	if passErr := utils.CheckPassword(req.Password, user.Password); passErr != nil {
		utils.RespondError(w, http.StatusUnauthorized, passErr, "invalid password")
		return
	}

	sessionID, sessionErr := dbHelper.CreateUserSession(user.ID)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create session")
		return
	}
	token, tokenErr := utils.GenerateJWT(user.ID, sessionID)
	if tokenErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, tokenErr, "failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	sessionID := userCtx.SessionID
	if delErr := dbHelper.DeleteUserSession(database.DB, sessionID); delErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, delErr, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID
	sessionID := userCtx.SessionID

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		delErr := dbHelper.DeleteUser(tx, userID)
		if delErr != nil {
			return delErr
		}
		return dbHelper.DeleteUserSession(tx, sessionID)
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to delete user account")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"account deleted successfully"})
}
