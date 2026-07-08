package middlewares

import (
	"My-todo-app/database/dbHelper"
	"My-todo-app/models"
	"My-todo-app/utils"
	"context"
	"net/http"
)

type ContextKeys string

const userContext ContextKeys = "userContext"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			http.Error(w, "invalid user id", http.StatusUnauthorized)
			return
		}

		sessionID, ok := claims["sessionID"].(string)
		if !ok {
			http.Error(w, "invalid session id", http.StatusUnauthorized)
			return
		}

		isValid, err := dbHelper.IsSessionValid(sessionID)

		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to validate session")
			return
		}

		if !isValid {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}
		userCtx := &models.UserCtx{
			UserID:    userID,
			SessionID: sessionID,
		}
		ctx := context.WithValue(r.Context(), userContext, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *models.UserCtx {
	user, ok := r.Context().Value(userContext).(*models.UserCtx)
	if !ok {
		return nil
	}
	return user
}
