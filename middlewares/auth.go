package middlewares

import (
	"My-todo-app/database/dbHelper"
	"My-todo-app/model"
	"My-todo-app/utils"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userId"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, err, "unauthorized")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userId"].(string)
		sessionID := claims["sessionId"].(string)

		isValid, err := dbHelper.IsSessionValid(sessionID)

		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to validate session")
			return
		}

		if !isValid {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}

		userCtx := model.UserCtx{
			UserID:    userID,
			SessionID: sessionID,
		}
		ctx := context.WithValue(r.Context(), UserKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
func UserContext(r *http.Request) model.UserCtx {
	return r.Context().Value(UserKey).(model.UserCtx)
}
