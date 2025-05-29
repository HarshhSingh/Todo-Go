package Routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main/handler"
	"net/http"
)

type publicRoutes struct {
	chi.Router
	server *http.Server
}

func PublicRoutes() *publicRoutes {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Get("/health-check", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("Hola ddd Amigo!"))
	})
	router.Get("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hola Amigo!"))
	})
	router.Route("/todo", func(todo chi.Router) {
		todo.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		todo.Get("/tasks", handler.GetTasks)
		todo.Post("/task", handler.PostTask)
		todo.Put("/task/{taskID}", handler.EditTask)
		todo.Delete("/task/{taskID}", handler.DeleteTask)

	})
	fmt.Printf("here chi router running %s \n", router)

	return &publicRoutes{
		Router: router,
	}
}
