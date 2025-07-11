package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/env"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/features/customer_sync"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/features/order"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func Run(ctx context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	serverErrors := make(chan error, 1)
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	e, err := env.ParseEnv()
	if err != nil {
		logFatal(err)
	}

	dbConn, err := newConnection(e.DB)
	if err != nil {
		logFatal(err)
	}
	defer func() { _ = dbConn.Close() }()

	streamEnv, err := newStreamEnvironment(e)
	if err != nil {
		logFatal(err)
	}
	defer func() { _ = streamEnv.Close() }()

	consumerUpdate := func(streamName string, isActive bool) stream.OffsetSpecification {
		// This function is called when the consumer is promoted to active
		// be careful with the logic here, it is called in the consumer thread
		// the code here should be fast, non-blocking and without side effects
		fmt.Printf("[%s] - Consumer promoted for: %s. Active status: %t\n", time.Now().Format(time.TimeOnly),
			streamName, isActive)

		// In this example, we store the offset server side and we retrieve it
		// when the consumer is promoted to active
		offset, err := streamEnv.QueryOffset(e.ConsumerName, streamName)
		if err != nil {
			// If the offset is not found, we start from the beginning
			return stream.OffsetSpecification{}.First()
		}

		// If the offset is found, we start from the last offset
		// we add 1 to the offset to start from the next message
		return stream.OffsetSpecification{}.Offset(offset + 1)
	}

	customerChan := make(chan customer_sync.CustomerEvent, 64)
	consumer, err := streamEnv.NewConsumer(
		api.StreamName,
		func(consumerContext stream.ConsumerContext, message *amqp.Message) {
			customer_sync.HandleMessage(dbConn, customerChan, message)
			err := consumerContext.Consumer.StoreOffset()
			if err != nil {
				log.Printf("Failed to store offset: %v", err)
			}
		},
		stream.NewConsumerOptions().
			SetConsumerName(e.ConsumerName).
			SetSingleActiveConsumer(stream.NewSingleActiveConsumer(consumerUpdate)),
	)
	if err != nil {
		logFatal(err)
	}
	defer func() { _ = consumer.Close() }()

	metrics.New()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	err = order.Setup(ctx, e, dbConn, customerChan, mux)
	if err != nil {
		logFatal(err)
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Addr:              e.HTTPServerAddress,
	}

	go func() {
		log.Printf("Starting order HTTP server on %s", e.HTTPServerAddress)
		serverErrors <- server.ListenAndServe()
	}()

	log.Printf("Order service consumer started for stream: %s", api.StreamName)

	select {
	case sig := <-sigChannel:
		log.Printf("Received signal: %s", sig.String())
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	case <-ctx.Done():
		log.Println("Context done")
	}

	log.Println("Shutting down order service...")
	cancelCtx()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}

func newConnection(config env.DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func newStreamEnvironment(e env.Env) (*stream.Environment, error) {
	streamEnv, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(e.RabbitMQHost).
			SetPort(e.RabbitMQPort).
			SetUser(e.RabbitMQUser).
			SetPassword(e.RabbitMQPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ environment: %w", err)
	}
	return streamEnv, nil
}

func logFatal(err error) {
	log.Printf("Fatal error: %v", err)
	os.Exit(1)
}
