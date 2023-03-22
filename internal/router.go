package togo

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Use this to allow specific origin hosts
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/signup", HandleUserSignup)
	r.Post("/login", HandleUserLogin)
	r.Get("/logout", HandleUserLogout)

	// Protected routes
	r.Route("/", func(r chi.Router) {
		r.Use(HandleAuthorizeRoute)

		r.Get("/me", HandleGetUserInfo)
		r.Put("/todos", HandleCreateTodo)
		r.Get("/todos", HandleGetTodoList)
	})

	return r
}
