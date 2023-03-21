package togo

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/signup", HandleUserSignup)
	r.Get("/login", HandleUserLogin)

	// Protected routes
	r.Route("/", func(r chi.Router) {
		r.Use(HandleAuthorizeRoute)

		r.Get("/me", HandleGetUserInfo)
		r.Put("/todos", HandleCreateTodo)
		r.Get("/todos", HandleGetTodoList)
	})

	return r
}
