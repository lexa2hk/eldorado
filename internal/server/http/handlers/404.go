package handlers

import (
	"net/http"

	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

func Handle404(w http.ResponseWriter, r *http.Request) error {
	return response.NotFound("resource")
}
