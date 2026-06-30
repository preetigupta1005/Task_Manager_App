package handler

import (
	"My-todo-app/database/dbHelper"
	"My-todo-app/middleware"
	"My-todo-app/model"
	"My-todo-app/utils"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userReq model.RegisterRequest
	err := utils.ParseBody(r, &userReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
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

	if delErr := dbHelper.DeleteUserSession(userCtx.SessionID); delErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, delErr, "failed to logout")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}
