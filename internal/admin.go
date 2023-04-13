package togo

import (
	"net/http"

	"go.uber.org/zap"
)

func HandleAuthorizeAdminRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, userIdOk := r.Context().Value(userIdContextKey{}).(int64)
		if !userIdOk {
			Log.Error(
				"Could not get user ID from context",
				zap.Any("username", r.Context().Value(userIdContextKey{})),
			)
			w.WriteHeader(http.StatusInternalServerError)
		}

		var user User
		err := Db.Get(&user, "SELECT * FROM users WHERE id = ?", userId)
		if err != nil {
			Log.Error("Could not get user from database", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if *user.Username != Config.AdminUsername {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HandleGetUserList(w http.ResponseWriter, r *http.Request) {
	var users []UserInfoAdmin
	err := Db.Select(&users, `
		SELECT u.id, u.username, u.created_at, count FROM users u INNER JOIN 
		( SELECT todos.owner, COUNT(*) count FROM todos GROUP BY todos.owner ) todo_count
		ON u.id = todo_count.owner;
	`)
	if err != nil {
		Log.Error("Could not get users from database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Render(w, r, Response{
		HTTPStatusCode: http.StatusOK,
		Data:           users,
		StatusText:     "success",
	})
}
