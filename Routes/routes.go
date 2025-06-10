package Routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main/handler"
	"main/middlewares"
	"net/http"
)

type publicRoutes struct {
	chi.Router
	server *http.Server
}

func PublicRoutes() *publicRoutes {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hola Amigo!"))
	})
	router.Post("/register", handler.RegisterUser)
	router.Post("/login", handler.LoginUser)

	router.Route("/todo", func(todo chi.Router) {
		todo.Use(middlewares.JWTAuthorisation)
		todo.Route("/", ProtectedRoutes)
	})
	return &publicRoutes{
		Router: router,
	}
}

func ProtectedRoutes(todo chi.Router) {
	fmt.Println("Protected Routes")
	todo.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	todo.Get("/{userId}/tasks", handler.AllTasksById)
	todo.Post("/{userId}/task", handler.PostTask)
	todo.Put("/{userId}/task/{taskID}", handler.EditTask)
	todo.Delete("/{userId}/task/{taskID}", handler.DeleteTask)
}
