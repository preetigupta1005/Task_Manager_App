package server

import (
	"context"
	"net/http"
	"time"

	"My-todo-app/handlers"
	"My-todo-app/middlewares"

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
	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/register", handlers.RegisterUser)
		v1.Post("/login", handlers.LoginUser)

		v1.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			r.Route("/user", func(user chi.Router) {
				user.Get("/profile", handlers.GetUser)
				user.Post("/logout", handlers.LogoutUser)
				user.Delete("/delete", handlers.DeleteUser)
			})
			r.Route("/todo", func(r chi.Router) {

				r.Get("/", handlers.GetAllTodos)
				r.Post("/", handlers.CreateTodo)
				r.Delete("/delete-all", handlers.DeleteAllTodos)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetTodoById)
					r.Put("/", handlers.UpdateTodo)
					r.Delete("/", handlers.DeleteTodo)
					r.Put("/complete", handlers.MarkTodoAsCompleted)
				})
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
