package togo

import (
	"log"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Global configuration

var Db *sqlx.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig

func InitGlobalConfig(env string, d *sqlx.DB) {
	err := initLogger(env)
	if err != nil {
		log.Fatal(err)
	}

	initDb(d)
	initBasicAuthConfig()
}

func initDb(d *sqlx.DB) {
	Db = d
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}
