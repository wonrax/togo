package togo

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

type Response struct {
	HTTPStatusCode int          `json:"-"` // http response status code
	Cookie         *http.Cookie `json:"-"` // http response cookie

	Data       interface{} `json:"data,omitempty"`
	StatusText string      `json:"status"`          // user-level status message
	ErrorText  string      `json:"error,omitempty"` // application-level error message, for debugging
}

func Render(w http.ResponseWriter, r *http.Request, response Response) {
	if response.HTTPStatusCode == 0 {
		response.HTTPStatusCode = http.StatusOK
	}

	if response.Cookie != nil {
		http.SetCookie(w, response.Cookie)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPStatusCode)

	json.NewEncoder(w).Encode(response)
}

func Bind(r *http.Request, data interface{}) error {
	if r.Header.Get("Content-Type") == "application/json" {
		return json.NewDecoder(r.Body).Decode(data)
	}
	return ErrUnsupportedContentType
}
