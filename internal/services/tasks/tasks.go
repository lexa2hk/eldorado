package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/cache"
	"github.com/romankravchuk/eldorado/internal/storages/cache/redis"
	"github.com/romankravchuk/eldorado/internal/storages/tasks"
	"github.com/romankravchuk/eldorado/internal/storages/tasks/pg"
)

type Option func(*Service) error

func WithTaskStorage(tasks tasks.Storage) Option {
	return func(s *Service) error {
		s.tasks = tasks
		return nil
	}
}

func WithTaskPostgresStorage(url string) Option {
	return func(s *Service) error {
		conn, err := storages.NewDBPool("postgres", url)
		if err != nil {
			return err
		}

		tasks, err := pg.New(conn)
		if err != nil {
			return err
		}

		return WithTaskStorage(tasks)(s)
	}
}

func WithCache(cache cache.Cache, ttl time.Duration) Option {
	return func(s *Service) error {
		s.cache = cache
		s.cacheTTL = ttl
		return nil
	}
}

func WithRedisCache(url string, ttl time.Duration) Option {
	return func(s *Service) error {
		client, err := storages.NewRedisClient(url)
		if err != nil {
			return err
		}

		return WithCache(redis.New(client), ttl)(s)
	}
}

type Service struct {
	tasks tasks.Storage

	cache    cache.Cache
	cacheTTL time.Duration
}

func New(opts ...Option) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Service) List(ctx context.Context, userID string) ([]data.Task, error) {
	cache, found, err := s.cache.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if found {
		tasks := make([]data.Task, 0)
		if err := json.Unmarshal(cache, &tasks); err != nil {
			return nil, err
		}
		return tasks, nil
	}

	tasks, err := s.tasks.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(tasks)
	if err := s.cache.Set(ctx, userID, data, s.cacheTTL); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) Create(ctx context.Context, userID string, t data.Task) (data.Task, error) {
	t.UserID = userID

	if err := s.tasks.Save(ctx, &t); err != nil {
		return data.Task{}, err
	}

	if err := s.cache.Del(ctx, userID); err != nil {
		return data.Task{}, err
	}

	return t, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.tasks.Delete(ctx, id); err != nil {
		return err
	}

	userID, ok := ctx.Value(api.UserIDKey).(string)
	if ok {
		if err := s.cache.Del(ctx, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id string, t data.Task) (data.Task, error) {
	t.ID = id

	if err := s.tasks.Update(ctx, &t); err != nil {
		return data.Task{}, err
	}

	userID, ok := ctx.Value(api.UserIDKey).(string)
	if ok {
		if err := s.cache.Del(ctx, userID); err != nil {
			return data.Task{}, err
		}
	}

	return t, nil
}
