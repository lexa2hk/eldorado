package response

import (
	"net/http"
)

type APIError struct {
	Status  int
	Message string
}

func (e APIError) Error() string {
	return e.Message
}

func NotFound(r string) APIError {
	return APIError{
		Status:  http.StatusNotFound,
		Message: r + " not found",
	}
}
