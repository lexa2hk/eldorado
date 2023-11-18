package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/romankravchuk/eldorado/internal/config"
	"github.com/romankravchuk/eldorado/internal/pkg/logger"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/services/auth"
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
	"google.golang.org/grpc"
)

var serviceRules = map[string]map[string]string{
	"sign-up": {
		"Email":    "required,email",
		"Username": "required,alpha,gte=5,lte=20",
		"Password": "required,gte=8,alphanum,lte=20",
	},
	"token": {
		"Email":    "required,email",
		"Password": "required,gte=8,alphanum,lte=20",
	},
	"refresh": {
		"Refresh": "required",
	},
}

func main() {
	cfg, err := config.LoadAuthServiceConfig()
	failOnError("failed to load config", err)

	log := logger.New(cfg.Env, os.Stderr)

	log.Debug("the debug mode is activated")
	log.Debug("config loaded", slog.Any("cfg", cfg))

	validator.RegisterRules(&proto.SignUpRequest{}, serviceRules["sign-up"])
	validator.RegisterRules(&proto.TokenRequest{}, serviceRules["token"])
	validator.RegisterRules(&proto.RefreshRequest{}, serviceRules["refresh"])

	svc, err := auth.New(
		auth.WithLogger(log),
		auth.WithUsersPostgresStorage(cfg.Postgres.URL),
		auth.WtihRedisSessionsStorage(cfg.Redis.URL),
		auth.WithAccessCreds(cfg.AccessCreds.PrivateKey, cfg.AccessCreds.PublicKey, cfg.AccessCreds.Expires),
		auth.WithRefreshCreds(cfg.RefreshCreds.PrivateKey, cfg.RefreshCreds.PublicKey, cfg.RefreshCreds.Expires),
	)
	failOnError("failed to create auth service", err)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	failOnError("failed to create listener", err)

	log.Info("the auth service started", slog.String("address", lis.Addr().String()))

	gsrv := grpc.NewServer()

	proto.RegisterAuthServiceServer(gsrv, svc)

	go func() {
		failOnError("failed to start auth service", gsrv.Serve(lis))
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	log.Info("the auth service stopped")

	gsrv.GracefulStop()
	os.Exit(0)
}

func failOnError(msg string, err error) {
	if err != nil {
		slog.Error(msg, sl.Err(err))
		os.Exit(1)
	}
}
