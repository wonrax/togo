package togo

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

var wrongAuthHeaderFormat = Response{
	HTTPStatusCode: http.StatusBadRequest,
	StatusText:     "Invalid Authorization header. Format is 'Basic <base64 encoded username:password>'.",
}

type userIdContextKey struct{}

func HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request and unmarshal it into a UserSignupRequest
	// struct
	var userSignupRequest UserSignupRequest
	err := Bind(r, &userSignupRequest)
	if err != nil {
		if err == ErrUnsupportedContentType {
			// TODO move this to a middleware
			Log.Debug("Unsupported content type", zap.String("content-type", r.Header.Get("Content-Type")))
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		Log.Error("Could not read request body", zap.Error(err))
		Render(w, r, Response{
			HTTPStatusCode: http.StatusBadRequest,
			StatusText:     "Could not read request body",
		})
		return
	}

	Log.Debug("Registering a user", zap.Any("credentials", userSignupRequest))

	// Check if the user already exists
	user, err := DbFind(
		map[string]interface{}{"username": *userSignupRequest.Username},
		[]string{"username"},
		"users",
	)
	if err != nil {
		Log.Error("Something went wrong", zap.Error(err))
		Render(w, r, Response{
			HTTPStatusCode: http.StatusInternalServerError,
			StatusText:     "We could not process your request",
		})
	}
	if len(user) > 0 {
		Log.Debug("User already exists", zap.Any("user", user))
		Render(w, r, Response{
			HTTPStatusCode: http.StatusConflict,
			StatusText:     "User already exists",
		})
		return
	}

	// generate iso 8601 timestamp for created_at and updated_at
	timestamp := time.Now().UTC().Format(time.RFC3339)

	key, salt := basicAuthConfig.HashPassword(*userSignupRequest.Password)
	_, err = DbInsert(
		map[string]interface{}{
			"username":        userSignupRequest.Username,
			"hashed_password": key,
			"password_salt":   salt,
			"created_at":      timestamp,
			"updated_at":      timestamp,
		},
		"users",
	)

	if err != nil {
		Log.Error("Could not register user", zap.Error(err))
		Render(w, r, Response{
			HTTPStatusCode: http.StatusBadRequest,
			StatusText:     "Could not register user",
		})
		return
	}

	// generate base64 encoded token
	token := base64.StdEncoding.EncodeToString(
		[]byte(*userSignupRequest.Username + ":" + *userSignupRequest.Password))

	Render(w, r, Response{
		HTTPStatusCode: http.StatusCreated,
		Data:           map[string]string{"token": token},
		StatusText:     "User created",
	})
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest UserLoginRequest
	err := Bind(r, &userLoginRequest)
	if err != nil {
		if err == ErrUnsupportedContentType {
			// TODO move this to a middleware
			Log.Debug("Unsupported content type", zap.String("content-type", r.Header.Get("Content-Type")))
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		Log.Error("Could not read request body", zap.Error(err))
		Render(w, r, Response{
			HTTPStatusCode: http.StatusBadRequest,
			StatusText:     "Could not read request body",
		})
		return
	}

	Log.Debug("Login requested", zap.Any("credentials", userLoginRequest))

	//  check if the user exists
	user, err := DbFind(
		map[string]interface{}{"username": *userLoginRequest.Username},
		[]string{"hashed_password", "password_salt"},
		"users",
	)
	if err != nil {
		Log.Error("Could not find user", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Log.Debug("Found login user", zap.Any("user", user))

	if len(user) == 0 {
		Render(w, r, Response{
			HTTPStatusCode: http.StatusNotFound,
			StatusText:     "User not found",
		})
		return
	}

	//  check if the password is correct
	if !basicAuthConfig.VerifyPassword(
		*userLoginRequest.Password,
		user[0]["hashed_password"].(string),
		user[0]["password_salt"].(string),
	) {
		Render(w, r, Response{
			HTTPStatusCode: http.StatusUnauthorized,
			StatusText:     "Invalid credentials",
		})
		return
	}

	// generate base64 encoded token
	token := base64.StdEncoding.EncodeToString(
		[]byte(*userLoginRequest.Username + ":" + *userLoginRequest.Password))

	Render(w, r, Response{
		HTTPStatusCode: http.StatusOK,
		Data:           map[string]string{"token": token},
		StatusText:     "Login successful",
	})
}

func HandleAuthorizeRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		Log.Debug("Authorization requested", zap.String("Authorization header", auth))

		authStr := strings.Split(auth, "Basic ")
		if len(authStr) != 2 {
			Log.Info("Invalid Authorization header")
			Render(w, r, wrongAuthHeaderFormat)
			return
		}

		rawDecodedText, err := base64.StdEncoding.DecodeString(authStr[1])
		if err != nil {
			Log.Info("Could not decode Authorization header", zap.Error(err))
			Render(w, r, wrongAuthHeaderFormat)
			return
		}

		Log.Debug("Decoded authorization header", zap.String("decoded text", string(rawDecodedText)))

		decodedText := strings.Split(string(rawDecodedText), ":")
		if len(decodedText) != 2 {
			Log.Info("Invalid Authorization header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := DbFind(
			map[string]interface{}{"username": decodedText[0]},
			[]string{"id", "hashed_password", "password_salt"},
			"users",
		)
		if err != nil {
			Log.Error("Could not find user", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		Log.Debug("Found login user", zap.Any("user", user))

		if len(user) == 0 {
			Render(w, r, Response{
				HTTPStatusCode: http.StatusNotFound,
				StatusText:     "User not found",
			})
			return
		}

		//  check if the password is correct
		if !basicAuthConfig.VerifyPassword(
			decodedText[1],
			user[0]["hashed_password"].(string),
			user[0]["password_salt"].(string),
		) {
			Render(w, r, Response{
				HTTPStatusCode: http.StatusUnauthorized,
				StatusText:     "Invalid credentials",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIdContextKey{}, user[0]["id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	userId, err := r.Context().Value(userIdContextKey{}).(int64)
	if !err {
		Log.Error(
			"Could not get user ID from context",
			zap.Any("username", r.Context().Value(userIdContextKey{})),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err_ := DbFind(
		map[string]interface{}{"id": userId},
		[]string{"username", "created_at"},
		"users",
	)
	if err_ != nil {
		Log.Error(
			"Something went wrong when looking for a user",
			zap.Int64("userId", userId),
			zap.Error(err_),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(user) == 0 {
		Log.Error(
			"Somehow this user does not exist anymore???",
			zap.Int64("userId", userId),
			zap.Error(err_),
		)
		w.WriteHeader(http.StatusBadRequest)
	}

	Render(w, r, Response{
		HTTPStatusCode: http.StatusOK,
		Data:           user[0],
		StatusText:     "success",
	})
}
