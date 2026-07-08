package handlers

import (
	"My-todo-app/database"
	"My-todo-app/database/dbHelper"
	"My-todo-app/middlewares"
	"My-todo-app/models"
	"My-todo-app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

var v = validator.New()

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userReq models.RegisterRequest
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

	hashedPassword, err := utils.HashedPassword(userReq.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to secure password")
		return
	}
	var token string
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(tx, userReq.Name, userReq.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}

		sessionID, sessErr := dbHelper.CreateUserSession(tx, userID)
		if sessErr != nil {
			return sessErr
		}
		var tokenErr error
		token, tokenErr = utils.GenerateJWT(userID, sessionID)
		return tokenErr
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to register user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{token})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

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

	sessionID, sessionErr := dbHelper.CreateUserSession(database.DB, user.ID)
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
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	if err := dbHelper.DeleteUserSession(database.DB, userCtx.SessionID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	userID := userCtx.UserID
	sessionID := userCtx.SessionID

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		if err := dbHelper.DeleteAllTodos(tx, userID); err != nil {
			return err
		}
		if err := dbHelper.DeleteUser(tx, userID); err != nil {
			return err
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	user, err := dbHelper.GetUser(userCtx.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get user")
		return
	}
	utils.RespondJSON(w, http.StatusOK, user)
}
