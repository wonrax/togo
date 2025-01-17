package togo

import (
	"fmt"
	"log"
	"os"
	"time"

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
	DatabasePath  string
}

var Db *sqlx.DB
var Log *zap.Logger
var basicAuthConfig *BasicAuthConfig
var Config *ConfigType
var Tables *TablesType

func InitGlobalConfig(env string) {
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
	Config.DatabasePath = os.Getenv("DATABASE_PATH")

	err = initLogger(env)
	if err != nil {
		log.Fatal(err)
	}

	// Database related
	initDb()
	Tables = &TablesType{
		Users: "users",
		Todos: "todos",
	}

	initBasicAuthConfig()

	// Init admin user
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminUsername != "" && adminPassword != "" {
		adminUser := &UserCredentials{
			Username: &adminUsername,
			Password: &adminPassword,
		}
		err = addAdminUser(adminUser)
		if err != nil {
			Log.Warn("Cannot update admin user", zap.Error(err))
		} else {
			Log.Info(fmt.Sprintf("Admin user `%s` updated", adminUsername))
			Config.AdminUsername = adminUsername
		}
	} else {
		Log.Warn("Admin user not configured")
	}
}

func initDb() {
	if Config.DatabasePath == "" {
		Config.DatabasePath = "togo.db"
		Log.Warn("Database path not set. Using default path: togo.db. The data could be lost every time the container is restarted. Set the DATABASE_PATH environment variable to a persistent path to avoid this.")
	}

	d, err := sqlx.Open("sqlite3", Config.DatabasePath)
	if err != nil {
		log.Fatal("Cannot open sqlite database", zap.Error(err))
	}

	Db = d

	// DATABASE MIGRATION
	DBMigrate()
}

func initBasicAuthConfig() {
	basicAuthConfig = NewBasicAuthConfig()
}

func addAdminUser(u *UserCredentials) error {
	adminUser := User{}
	err := Db.Get(&adminUser, "SELECT * FROM users WHERE username = ?", u.Username)

	key, salt := basicAuthConfig.HashPassword(*u.Password)
	// generate iso 8601 timestamp for updated_at
	timestamp := time.Now().UTC().Format(time.RFC3339)

	if err != nil {
		Log.Debug("Admin user does not exist. Creating one.")
		_, err := Db.Exec(`
			INSERT INTO users (username, hashed_password, password_salt, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`, u.Username, key, salt, timestamp, timestamp)

		if err != nil {
			return err
		}
	} else {
		// If exists, update the password
		Log.Debug("Admin user exists. Updating password.")
		_, err := Db.Exec(`
			UPDATE users SET hashed_password = ?, password_salt = ?, updated_at = ?
			WHERE username = ?
		`, key, salt, timestamp, u.Username)
		if err != nil {
			return err
		}
	}
	return nil
}
