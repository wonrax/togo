package togo

import (
	"database/sql"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global configuration

var Db *sql.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig

func InitGlobalConfig(d *sql.DB) {
	initDb(d)
	initLogger()
	initBasicAuthConfig()
}

func initDb(d *sql.DB) {
	Db = d
}

func initLogger() error {
	var logConfig = zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var err error
	Log, err = logConfig.Build()

	return err
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}
