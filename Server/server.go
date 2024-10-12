package Server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	Handle "todo/Handler"
	"todo/Middleware"
	"todo/Utils"
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

func SetupRoutes() *Server {

	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Get("/status", Utils.Health)
		r.Post("/register", Handle.CreateUser)
		r.Post("/login", Handle.UserLogin)
	})

	router.Group(func(r chi.Router) {
		r.Use(Middleware.Authentication)
		r.Post("/v1/deleteAccount", Handle.DeleteUser)
		r.Post("/v1/Logout", Handle.Logout)
		r.Route("/v1/todos", func(r chi.Router) {
			r.Post("/create", Handle.CreateNote)
			r.Get("/search", Handle.GetTodoByName)
			r.Put("/updateStatus", Handle.MarkCompleted)
			r.Delete("/deleteTodo", Handle.TodoDeleted)
		})
	})

	return &Server{
		Router: router,
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
