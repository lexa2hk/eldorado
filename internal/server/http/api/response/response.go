package response

import (
	"encoding/json"
	"net/http"
)

type M map[string]any

func JSON(w http.ResponseWriter, code int, data M) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
