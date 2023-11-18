package main

import (
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"log/slog"
	"os"
)

func main() {
	//cfg, err := config.LoadStatisticServiceConfig()
	//failOnError("failed to load config", err)
	//
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	//
	//log := logger.New(cfg.Env, os.Stderr)
	//
	//log.Debug("the debug mode is activated")
	//
	//svc, err := statistic.New(
	//	statistic.WithCron(cfg.Schedule),
	//	statistic.WithLogger(log),
	//	statistic.WithPosgresTasksStorage(cfg.Postgres.URL),
	//	statistic.WithRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName),
	//)
	//failOnError("failed to create statistic service", err)
	//
	//go func() {
	//	log.Info("the statistic service started")
	//	if err := svc.Start(); err != nil {
	//		log.Error("failed to start statistic service", sl.Err(err))
	//	}
	//}()
	//
	//<-sigCh
	//log.Info("the statistic service stopped")
}

func failOnError(msg string, err error) {
	if err != nil {
		slog.Error(msg, sl.Err(err))
		os.Exit(1)
	}
}
