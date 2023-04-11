package togo

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func Setup(environment string) {
	// DATABASE INITIALIZATION
	db, err := sqlx.Open("sqlite3", "togo.db")
	if err != nil {
		log.Fatal("Cannot open sqlite database", zap.Error(err))
	}

	InitGlobalConfig(environment, db)

	// DATABASE MIGRATION
	DBMigrate()
}

func Cleanup() {
	if Db != nil {
		Db.Close()
	}
}

func Start(environment string) {
	Setup(environment)
	defer Cleanup()
	// ROUTING
	port := "3000"
	Log.Info("Starting server on port " + port)
	http.ListenAndServe(":"+port, CreateRouter())
}
