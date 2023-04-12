package togo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
