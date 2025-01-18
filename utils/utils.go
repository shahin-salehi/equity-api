package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DeserializeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func SerializeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Ttpe", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	SerializeJSON(w, status, map[string]string{"error": err.Error()})
}
