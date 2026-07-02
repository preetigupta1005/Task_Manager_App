package middlewares

import (
	"My-todo-app/models"
	"context"
	"net/http"

	"My-todo-app/database/dbHelper"

	"github.com/sirupsen/logrus"
)

type ContextKeys string

const userContext ContextKeys = "userContext"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("x-api-key")
		if sessionID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := dbHelper.GetUserBySession(sessionID)
		if err != nil || user == nil {
			logrus.WithError(err).Errorf("failed to get user with session: %s", sessionID)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContext).(*models.User)
	if !ok {
		return nil
	}
	return user
}
