package togo

import (
	"encoding/base64"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

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

	key, salt := basicAuthConfig.HashPassword(*userSignupRequest.Password)
	_, err = DbInsert(
		map[string]interface{}{
			"username":        userSignupRequest.Username,
			"hashed_password": key,
			"password_salt":   salt,
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

var wrongAuthHeaderFormat = Response{
	HTTPStatusCode: http.StatusBadRequest,
	StatusText:     "Invalid Authorization header. Format is 'Basic <base64 encoded username:password>'.",
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	Log.Debug("Login requested", zap.String("Authorization header", auth))

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

	Render(w, r, Response{
		HTTPStatusCode: http.StatusOK,
		StatusText:     "Login successful",
	})
}
