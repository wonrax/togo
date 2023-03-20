package togo

import (
	"encoding/base64"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type UserSignupRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request and unmarshal it into a UserSignupRequest
	// struct
	var userSignupRequest UserSignupRequest
	err := ReadJSONBody(w, r, &userSignupRequest)
	if err != nil {
		Log.Error("Could not read request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Log.Debug("Registering a user", zap.Any("request", userSignupRequest))

	key, salt := basicAuthConfig.HashPassword(*userSignupRequest.Password)
	result, err := DbInsert(
		map[string]interface{}{
			"username":        userSignupRequest.Username,
			"hashed_password": key,
			"password_salt":   salt,
		},
		"users",
	)

	Log.Debug("Result of user registration", zap.Any("result", result))

	if err != nil {
		Log.Error("Could not register user", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	Log.Debug("Login requested", zap.String("Authorization header", auth))

	authStr := strings.Split(auth, "Basic ")
	if len(authStr) != 2 {
		Log.Info("Invalid Authorization header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(authStr[1])
	if err != nil {
		Log.Info("Could not decode Authorization header", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Log.Debug("Decoded authorization header", zap.String("decoded text", string(rawDecodedText)))

	decodedText := strings.Split(string(rawDecodedText), ":")
	if len(decodedText) != 2 {
		Log.Info("Invalid Authorization header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO look up the user in the database and compare the password
}
