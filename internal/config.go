package togo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Global configuration

var Db *sqlx.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig

func InitGlobalConfig(env string, d *sqlx.DB) {
	initDb(d)
	initLogger(env)
	initBasicAuthConfig()
}

func initDb(d *sqlx.DB) {
	Db = d
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}
