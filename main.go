package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	togo "github.com/wonrax/togo/src"
	"github.com/wonrax/togo/src/middleware"
	"go.uber.org/zap"
)

func UserSignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Signup page %s", r.URL.Path[1:])
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	// INIT LOGGER
	err := togo.InitLogger()
	if err != nil {
		log.Fatal("Could not init logger")
	}

	// DATABASE INITIALIZATION
	db, err := sql.Open("sqlite3", "togo.db")
	if err != nil {
		togo.Log.Fatal("Cannot open sqlite database", zap.Error(err))
	}
	defer db.Close()

	// DATABASE MIGRATION
	togo.DBMigrate(db)

	// ROUTING
	r := chi.NewRouter()
	r.Use(middleware.New(togo.Log, nil))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(id))
	})
	http.ListenAndServe(":3000", r)

}
