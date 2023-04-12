package togo

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Global configuration

type TablesType struct {
	Users string
	Todos string
}

type ConfigType struct {
	Environment   string
	AppURL        string
	AdminUsername string
}

var Db *sqlx.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig
var Config *ConfigType
var Tables *TablesType

func InitGlobalConfig(env string, d *sqlx.DB) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Expecting environment variables to be set.")
	}

	if env == "" {
		// if env is not overridden, use the ENVIRONMENT environment variable
		env = os.Getenv("ENVIRONMENT")
	}

	Config = &ConfigType{}
	Config.Environment = env
	Config.AppURL = os.Getenv("APP_URL")

	err = initLogger(env)
	if err != nil {
		log.Fatal(err)
	}

	// Database related
	initDb(d)
	Tables = &TablesType{
		Users: "users",
		Todos: "todos",
	}

	initBasicAuthConfig()
}

func initDb(d *sqlx.DB) {
	Db = d
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}
