package tasks

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name TaskUpdater
type TaskUpdater interface {
	Update(ctx context.Context, id string, t data.Task) (data.Task, error)
}

func HandleUpdateTask(log *slog.Logger, updater TaskUpdater) api.APIFunc {
	const op = "server.http.handlers.tasks.UpdateTask"

	type req struct {
		Title       string `json:"title" validate:"required,min=3,max=100"`
		Description string `json:"description" validate:"required,min=3,max=255"`
		IsCompleted bool   `json:"is_completed" validate:"boolean"`
	}

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

		input := new(req)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: msg,
			}
		}

		if err := validator.ValidateStruct(*input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 150*time.Millisecond)
		defer cancel()

		updated, err := updater.Update(ctx, chi.URLParam(r, "id"), data.Task{
			UserID:      userID,
			Title:       input.Title,
			Description: input.Description,
			IsCompleted: input.IsCompleted,
		})
		if err != nil {
			msg := "internal server error"

			log.Error(msg,
				sl.Err(err),
				slog.String("user_id", userID),
				slog.String("task_id", chi.URLParam(r, "id")),
				slog.Any("request body", input),
			)

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{
			"task": task{
				ID:          updated.ID,
				Title:       updated.Title,
				Description: updated.Description,
				CreatedOn:   updated.CreatedOn.Format(time.RFC3339),
				IsCompleted: updated.IsCompleted,
			},
		})
	}
}
