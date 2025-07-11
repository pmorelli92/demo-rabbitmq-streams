package customer_sync

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/env"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func Setup(
	e env.Env,
	db *pgxpool.Pool,
	logger *slog.Logger,
	streamEnv *stream.Environment) error {

	handler := handler{
		db:     db,
		logger: logger,
	}

	// This could be in a lib
	consumerUpdate := func(streamName string, isActive bool) stream.OffsetSpecification {
		logger.Info("Consumer promoted")
		offset, err := streamEnv.QueryOffset(e.RabbitMQ.ConsumerName, streamName)
		if err != nil {
			return stream.OffsetSpecification{}.First()
		}
		return stream.OffsetSpecification{}.Offset(offset + 1)
	}

	// This could be in a lib as well
	_, err := streamEnv.NewConsumer(
		api.StreamName,
		func(consumerContext stream.ConsumerContext, message *amqp.Message) {
			ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
			err := handler.consume(ctx, message)
			if err != nil {
				logger.Error("failed to consume message", "error", err)
			}
			cancelCtx()
			err = consumerContext.Consumer.StoreOffset()
			if err != nil {
				logger.Error("failed to store offset", "error", err)
			}
		},
		stream.NewConsumerOptions().
			SetConsumerName(e.RabbitMQ.ConsumerName).
			SetSingleActiveConsumer(stream.NewSingleActiveConsumer(consumerUpdate)),
	)

	return err
}
