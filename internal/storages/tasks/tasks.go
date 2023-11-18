package tasks

import (
	"context"
	"errors"

	"github.com/romankravchuk/eldorado/internal/data"
)

var ErrNotFound = errors.New("the task not found")

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Storage
type Storage interface {
	FindByUserID(ctx context.Context, userID string) ([]data.Task, error)
	UncompletedStatistic(ctx context.Context) ([]data.StatisticTask, error)
	Save(ctx context.Context, task *data.Task) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task *data.Task) error
}
