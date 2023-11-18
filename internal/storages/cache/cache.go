package cache

import (
	"context"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Cache
type Cache interface {
	Set(context.Context, string, []byte, time.Duration) error
	Get(context.Context, string) ([]byte, bool, error)
	Del(context.Context, string) error
}
