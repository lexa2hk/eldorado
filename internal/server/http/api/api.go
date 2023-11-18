package api

import (
	"context"
	"net/http"

	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHTTPHandlerFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			switch err := err.(type) {
			case response.APIError:
				response.JSON(w, err.Status, response.M{"error": err.Message})
			default:
				response.JSON(w, http.StatusInternalServerError, response.M{"error": err.Error()})
			}
		}
	}
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
