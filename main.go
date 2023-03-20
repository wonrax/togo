package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	togo "github.com/wonrax/togo/internal"
	"go.uber.org/zap"
)

func UserSignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Signup page %s", r.URL.Path[1:])
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	// DATABASE INITIALIZATION
	db, err := sql.Open("sqlite3", "togo.db")
	if err != nil {
		log.Fatal("Cannot open sqlite database", zap.Error(err))
	}
	defer db.Close()

	togo.InitGlobalConfig(db)

	// DATABASE MIGRATION
	togo.DBMigrate()

	// ROUTING
	togo.Log.Info("Starting server on port 3000")
	http.ListenAndServe(":3000", togo.CreateRouter())
}
