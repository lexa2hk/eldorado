package handlers

import (
	"net/http"

	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	return response.JSON(w, http.StatusOK, response.M{"message": "alive"})
}
