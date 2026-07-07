package middlewares

import (
	"My-todo-app/utils"
	"context"
	"net/http"
	"strings"
)

type ContextKeys string

const userContext ContextKeys = "userContext"

type UserCtx struct {
	UserID string
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userCtx := &UserCtx{
			UserID: claims["userID"].(string),
		}

		ctx := context.WithValue(r.Context(), userContext, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *UserCtx {
	user, ok := r.Context().Value(userContext).(*UserCtx)
	if !ok {
		return nil
	}
	return user
}
