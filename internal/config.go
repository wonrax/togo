package togo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
