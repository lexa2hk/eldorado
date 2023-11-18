package tasks

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name TasksLister
type TasksLister interface {
	List(ctx context.Context, userID string) ([]data.Task, error)
}

func HandleGetTasks(log *slog.Logger, lister TasksLister) api.APIFunc {
	const op = "server.http.handlers.tasks.GetTasks"

	type task struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedOn   string `json:"created_at"`
		IsCompleted bool   `json:"is_completed"`
	}
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

		tasks, err := lister.List(ctx, userID)
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err), slog.String("user_id", userID))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}

		objs := make([]task, len(tasks))
		for i, t := range tasks {
			objs[i] = task{
				ID:          t.ID,
				Title:       t.Title,
				Description: t.Description,
				CreatedOn:   t.CreatedOn.Format(time.RFC3339),
				IsCompleted: t.IsCompleted,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{
			"tasks": objs,
		})
	}
}
