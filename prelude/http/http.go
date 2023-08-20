package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Respond(w http.ResponseWriter, r *http.Request, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Could not encode in json", code)
		}
	}
}

func Created(w http.ResponseWriter, r *http.Request, id string) {
	path := r.URL.Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	w.Header().Add("Location", "//"+r.Host+path+id)
	Respond(w, r, nil, http.StatusCreated)
}

func Decode(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {
	return json.NewDecoder(r.Body).Decode(data)
}

func Echo(message string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) { Respond(rw, r, message, http.StatusOK) }
}

func FileServer(prefix, dirname string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(dirname)))
}
