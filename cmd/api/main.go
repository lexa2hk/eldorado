package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/config"
	"github.com/romankravchuk/eldorado/internal/pkg/logger"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/handlers"
	authhandlers "github.com/romankravchuk/eldorado/internal/server/http/handlers/auth"
	taskshandlers "github.com/romankravchuk/eldorado/internal/server/http/handlers/tasks"
	"github.com/romankravchuk/eldorado/internal/server/http/middleware"
	"github.com/romankravchuk/eldorado/internal/services/auth/client"
	"github.com/romankravchuk/eldorado/internal/services/tasks"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func main() {
	cfg, err := config.LoadApiConfig()
	if err != nil {
		slog.Error("failed to load configuration for api", sl.Err(err))
		os.Exit(1)
	}

	log := logger.New(cfg.Env, os.Stderr)

	authClient, err := client.New(cfg.AuthServiceAddr)
	if err != nil {
		slog.Error("failed to create auth service client", sl.Err(err))
		os.Exit(1)
	}

	svc, err := tasks.New(
		tasks.WithTaskPostgresStorage(cfg.Postgres.URL),
		tasks.WithRedisCache(cfg.Redis.URL, cfg.Redis.TTL),
	)
	if err != nil {
		slog.Error("failed to create tasks service", sl.Err(err))
		os.Exit(1)
	}

	mux := chi.NewMux()
	mux.NotFound(api.MakeHTTPHandlerFunc(handlers.Handle404))
	mux.MethodNotAllowed(api.MakeHTTPHandlerFunc(handlers.Handle404))

	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.Logger(log))
	mux.Use(chimiddleware.Recoverer)

	mux.Get("/health", api.MakeHTTPHandlerFunc(handlers.HandleHealthCheck))
	mux.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/", api.MakeHTTPHandlerFunc(authhandlers.HandleRegister(log, authClient)))
			r.Post("/token", api.MakeHTTPHandlerFunc(authhandlers.HandleGetTokenPairs(log, authClient)))
			r.Post("/refresh", api.MakeHTTPHandlerFunc(authhandlers.HandleRefreshToken(log, authClient)))
		})
		r.With(middleware.JWT(log, authClient)).Route("/tasks", func(r chi.Router) {
			r.Post("/", api.MakeHTTPHandlerFunc(taskshandlers.HandleCreateTask(log, svc)))
			r.Get("/", api.MakeHTTPHandlerFunc(taskshandlers.HandleGetTasks(log, svc)))
			r.Route("/{id}", func(r chi.Router) {
				r.Put("/", api.MakeHTTPHandlerFunc(taskshandlers.HandleUpdateTask(log, svc)))
				r.Delete("/", api.MakeHTTPHandlerFunc(taskshandlers.HandleDeleteTask(log, svc)))
			})
		})
	})

	srv := http.Server{
		Handler:      mux,
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("server starting", slog.String("addr", cfg.Server.Addr))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("failed to serve server", sl.Err(err))
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("failed to gracefully shutdown server", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("server successfully sutdown")
	os.Exit(0)
}
