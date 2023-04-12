package togo

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func CreateRouter() *chi.Mux {
	appURL := Config.AppURL
	if appURL == "" {
		appURL = "http://localhost:8088"
		Log.Warn("APP_URL environment variable not set. Defaulting to " + appURL)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{appURL}, // Use this to allow specific origin hosts
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
		r.Delete("/todos/{id}", HandleDeleteTodo)

		r.Route("/admin", func(r chi.Router) {
			r.Use(HandleAuthorizeAdminRoute)
			r.Get("/users", HandleGetUserList)
			// r.Get("/users/{id}", HandleGetUserInfo) TODO
			// r.Put("/users/{id}", HandleUpdateUser) TODO
			// r.Delete("/users/{id}", HandleDeleteUser) TODO
		})
	})

	return r
}
