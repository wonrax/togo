package togo

import (
	"fmt"

	"go.uber.org/zap"
)

type migration struct {
	name string
	sql  string
}

func DBMigrate() {
	migrations := []migration{
		{
			name: "init db",
			sql: `
				CREATE TABLE IF NOT EXISTS users(
					id INTEGER PRIMARY KEY,
					username TEXT NOT NULL,
					hashed_password TEXT NOT NULL,
					password_salt TEXT NOT NULL,
					created_at TEXT NOT NULL,
					updated_at TEXT NOT NULL
				);
				CREATE UNIQUE INDEX IF NOT EXISTS
					idx_users_username ON users (username);
			`,
		},
		{
			name: "create todo table",
			sql: `
				CREATE TABLE IF NOT EXISTS todos(
					id INTEGER PRIMARY KEY,
					owner INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
					title TEXT NULL,
					description TEXT NULL,
					completed INT NOT NULL DEFAULT 0,
					created_at TEXT NOT NULL,
					updated_at TEXT NOT NULL
				);
				CREATE INDEX IF NOT EXISTS
					idx_todos_updated_at ON todos (updated_at);
			`,
		},
	}

	for i, m := range migrations {
		result, err := Db.Exec(m.sql)
		if err != nil {
			Log.Panic("Database migration failed, stopping...", zap.Error(err))
		}

		rowsAffected, _ := result.RowsAffected()
		lastInsertId, _ := result.LastInsertId()

		Log.Info(
			fmt.Sprintf("Run migration #%d %s", i, m.name),
			zap.Int64("rows affected", rowsAffected),
			zap.Int64("last insert id", lastInsertId),
		)
	}
}
