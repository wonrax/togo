package togo_test

import (
	"testing"

	togo "github.com/wonrax/togo/internal"
)

func BenchmarkDbInsertPerformance(b *testing.B) {
	b.StopTimer()

	togo.Setup("test")
	defer togo.Cleanup()

	// Delete user with username test_user
	togo.Db.Exec("DELETE FROM users WHERE username = ?", "test_user")

	r, err := togo.DbInsert(
		map[string]interface{}{
			"username":        "test_user",
			"hashed_password": "hashed_password",
			"password_salt":   "password_salt",
			"created_at":      "timestamp",
			"updated_at":      "timestamp",
		},
		"users",
	)
	if err != nil {
		b.Fatal(err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		b.Fatal(err)
	}
	todo := map[string]interface{}{
		"owner":       id,
		"title":       "_test_title",
		"description": "_test_description",
		"completed":   false,
		"created_at":  "timestamp",
		"updated_at":  "timestamp",
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		togo.DbInsert(
			todo,
			"todos",
		)
	}
}
