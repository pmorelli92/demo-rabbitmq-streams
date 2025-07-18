package customer_sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	customerapi "github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	gen_sql "github.com/pmorelli92/demo-rabbitmq-streams/order/database/generated"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/metrics"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
)

type handler struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func (h handler) consume(ctx context.Context, message *amqp.Message) error {
	eventType, ok := message.Annotations["event_type"]
	if !ok {
		metrics.CustomerEventErrors.Inc()
		return fmt.Errorf("event_type annotation not found in message")
	}

	switch eventType {
	case customerapi.CustomerCreated:
		evt := customerapi.CustomerCreatedEvent{}
		if err := json.Unmarshal(message.Data[0], &evt); err != nil {
			metrics.CustomerEventErrors.Inc()
			return err
		}
		return h.update(ctx, evt.CustomerID, evt.Address, evt.Timestamp)

	case customerapi.CustomerAddressUpdated:
		evt := customerapi.CustomerAddressUpdatedEvent{}
		if err := json.Unmarshal(message.Data[0], &evt); err != nil {
			metrics.CustomerEventErrors.Inc()
			return err
		}
		return h.update(ctx, evt.CustomerID, evt.Address, evt.Timestamp)

	default:
		metrics.CustomerEventErrors.Inc()
		return fmt.Errorf("unknown event type: %s", eventType)
	}
}

func (h handler) update(
	ctx context.Context,
	customerID, address string,
	updatedAt time.Time) error {

	q := gen_sql.New(h.db)
	err := q.UpsertCustomer(ctx, gen_sql.UpsertCustomerParams{
		ID:      customerID,
		Address: address,
		UpdatedAt: pgtype.Timestamptz{
			Time:  updatedAt,
			Valid: true,
		},
	})

	if err != nil {
		metrics.CustomerEventErrors.Inc()
		return err
	}

	metrics.CustomersProcessed.Inc()
	h.logger.Info("customer updated", "customerId", customerID)
	return nil
}
