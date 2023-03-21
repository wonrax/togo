package togo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Global configuration

var Db *sqlx.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig

func InitGlobalConfig(d *sqlx.DB) {
	initDb(d)
	initLogger()
	initBasicAuthConfig()
}

func initDb(d *sqlx.DB) {
	Db = d
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}
