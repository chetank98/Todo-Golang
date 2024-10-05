package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"todo/Database"
	CreateTodo "todo/Handler"
)

// try to keep main.go in cmd folder
// you can write this single line message directly and no need to use w.write
// create an anonymous struct and encode it

func Home(w http.ResponseWriter, r *http.Request) {

	// handle error

	w.Write([]byte("Welcome to the app - Home Page of the todo "))
}

func main() {

	fmt.Println("Welcome to the app - ToDo")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// login route should be there
	// try to group the routes in multiple categories such as public routes and protected routes
	// routes that does not need auth keep them in public route
	// routes that needs auth keep it in private routes

	/*
		1) try to keep all routes in separate file server.go
		2) add function for server shutdown
	*/

	r.Route("/v1", func(r chi.Router) {
		r.Get("/home", Home)
		r.Get("/allNotes", CreateTodo.GetAll)
		r.Get("/Notes/{id}", CreateTodo.GetById)
		r.Post("/Note", CreateTodo.CreateNote)
		r.Put("/updateNote", CreateTodo.UpdateTodo)
		r.Delete("/Note/{id}", CreateTodo.DeleteNote)

		//user
		r.Post("/createuser", CreateTodo.CreateUser)

	})

	// Handle error while running the server

	http.ListenAndServe(":8000", r)

	err := Database.DBConnection
	if err != nil {
		log.Printf("error connecting DaB %v", err)
	}

	// handle error

	defer Database.CloseDatabase()

}
