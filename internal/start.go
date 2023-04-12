package togo

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Seperate starting the server and setting up the app
// because we want to be able to test the app without
// starting the server

func Setup(environment string) {
	InitGlobalConfig(environment)
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
