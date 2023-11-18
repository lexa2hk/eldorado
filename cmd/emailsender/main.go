package main

import (
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"log/slog"
	"os"
)

func main() {
	//cfg, err := config.LoadEmailSenderConfig()
	//failOnError("failed to load config", err)
	//
	//log := logger.New(cfg.Env, os.Stderr)
	//
	//log.Debug("debug mode is activated")
	//
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//
	//svc, err := emailsender.New(
	//	emailsender.WithLogger(log),
	//	emailsender.WithRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName),
	//	emailsender.WithSMTP(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Email, cfg.SMTP.Password),
	//)
	//failOnError("failed to created email sender service", err)
	//
	//go func() {
	//	log.Info("the email service started")
	//	if err := svc.Start(); err != nil {
	//		log.Error("failed to start email sender service", sl.Err(err))
	//	}
	//}()
	//
	//<-sigCh
	//svc.Close()
	//log.Info("the email service stopped")
}

func failOnError(msg string, err error) {
	if err != nil {
		slog.Error(msg, sl.Err(err))
		os.Exit(1)
	}
}
