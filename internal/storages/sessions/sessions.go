package sessions

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("the value in sessions storage not found")

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Storage
type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
