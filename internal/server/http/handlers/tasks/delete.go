package tasks

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name TaskDeleter
type TaskDeleter interface {
	Delete(ctx context.Context, id string) error
}

func HandleDeleteTask(log *slog.Logger, deleter TaskDeleter) api.APIFunc {
	const op = "server.http.handlers.tasks.DeleteTask"

	return func(w http.ResponseWriter, r *http.Request) error {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, ok := r.Context().Value(api.UserIDKey).(string)
		if !ok {
			msg := "forbidden"

			log.Error(msg, slog.String("error", "no user id in context"))

			return response.APIError{
				Status:  http.StatusForbidden,
				Message: msg,
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 150*time.Millisecond)
		defer cancel()

		if err := deleter.Delete(ctx, chi.URLParam(r, "id")); err != nil {
			msg := "internal server error"

			log.Error(msg,
				sl.Err(err),
				slog.String("user_id", userID),
				slog.String("task_id", chi.URLParam(r, "id")),
			)

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{"message": "ok"})
	}
}
