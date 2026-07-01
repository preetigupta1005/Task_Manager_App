package server

import (
	"context"
	"net/http"
	"time"

	"My-todo-app/handler"
	"My-todo-app/middleware"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetUpRoutes() *Server {
	r := chi.NewRouter()
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Route("/user", func(user chi.Router) {
			user.Get("/profile", handler.GetUser)
			user.Post("/logout", handler.LogoutUser)
			user.Delete("/delete", handler.DeleteUser)
		})
		r.Route("/todo", func(r chi.Router) {
			r.Get("/", handler.GetAllTodos)
			r.Post("/", handler.CreateTodo)
			r.Delete("/delete-all", handler.DeleteAllTodos)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetTodoById)
				r.Put("/", handler.UpdateTodo)
				r.Delete("/", handler.DeleteTodo)
				r.Put("/complete", handler.MarkTodoAsCompleted)
			})
		})
	})
	return &Server{
		Router: r,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
