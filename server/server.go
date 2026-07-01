package server

import (
	"My-todo-app/handler"
	"My-todo-app/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetUpRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Route("/user", func(user chi.Router) {
			//user.Get("/profile", handlers.GetUser)
			user.Post("/logout", handler.LogoutUser)
			user.Delete("/delete", handler.DeleteUser)
		})

		r.Route("/todo", func(r chi.Router) {
			r.Get("/", handler.GetAllTodos)
			r.Post("/", handler.CreateTodo)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetTodoById)
				r.Put("/", handler.UpdateTodo)
				r.Delete("/", handler.DeleteTodo)
				r.Put("/complete", handler.MarkTodoAsCompleted)
			})
		})
	})
	return r
}
