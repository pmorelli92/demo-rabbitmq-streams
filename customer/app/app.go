package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/database"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/env"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/features/customer"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func Run(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancelCtx := context.WithCancel(ctx)

	serverErrors := make(chan error, 1)
	defer close(serverErrors)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigChannel)

	e, err := env.ParseEnv()
	if err != nil {
		logFatal(logger, err)
	}

	if err = database.Migrate(e.DBConnectionMigrate); err != nil {
		logFatal(logger, err)
	}

	db, err := pgxpool.New(ctx, e.DBConnectionDSN)
	if err != nil {
		logFatal(logger, err)
	}

	streamEnv, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(e.RabbitMQ.Host).
			SetPort(e.RabbitMQ.Port).
			SetUser(e.RabbitMQ.User).
			SetPassword(e.RabbitMQ.Password),
	)
	if err != nil {
		logFatal(logger, err)
	}

	err = streamEnv.DeclareStream(api.StreamName, &stream.StreamOptions{
		MaxLengthBytes: stream.ByteCapacity{}.GB(1),
	})
	if err != nil && !errors.Is(err, stream.StreamAlreadyExists) {
		logFatal(logger, err)
	}

	metrics.New()
	mux := http.NewServeMux()

	// Features
	{
		err = customer.Setup(db, logger, streamEnv, mux)
		if err != nil {
			logFatal(logger, err)
		}
	}

	mux.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Handler:           http.TimeoutHandler(mux, e.HTTPTimeout, "request timed out"),
		Addr:              e.HTTPAddress,
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		logger.Info("starting HTTP server", "address", e.HTTPAddress)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case sig := <-sigChannel:
		logger.Info("received signal", "signal", sig.String())

	case err := <-serverErrors:
		logger.Error("server error", "error", err.Error())

	case <-ctx.Done():
		logger.Info("context done")
	}

	cancelCtx()
	logger.Info("shutting down...")

	db.Close()
	_ = httpServer.Close()
	_ = streamEnv.Close()
}

func logFatal(logger *slog.Logger, err error) {
	logger.Error(err.Error())
	os.Exit(1)
}
