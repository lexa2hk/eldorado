package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/eldorado/internal/storages/sessions"
)

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := s.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, sessions.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
