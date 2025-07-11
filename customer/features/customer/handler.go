package customer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	gen_sql "github.com/pmorelli92/demo-rabbitmq-streams/customer/database/generated"
	"github.com/pmorelli92/demo-rabbitmq-streams/customer/metrics"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/ha"
)

type handler struct {
	db       *pgxpool.Pool
	logger   *slog.Logger
	producer *ha.ReliableProducer
}

func (h *handler) createCustomer(ctx context.Context, rq createCustomerRq) (customerRs, error) {
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return customerRs{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := gen_sql.New(tx)
	customerID := uuid.NewString()

	err = q.InsertCustomer(ctx, gen_sql.InsertCustomerParams{
		ID:      customerID,
		Name:    rq.Name,
		Email:   rq.Email,
		Address: rq.Address,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return customerRs{}, err
	}

	err = h.publishEvent(
		uuid.NewString(),
		api.CustomerAddressUpdated,
		api.CustomerCreatedEvent{
			CustomerID: customerID,
			Timestamp:  time.Now(),
			Name:       rq.Name,
			Email:      rq.Email,
			Address:    rq.Address,
		})
	if err != nil {
		return customerRs{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return customerRs{}, err
	}

	return customerRs{
		CustomerID: customerID,
		Name:       rq.Name,
		Email:      rq.Email,
		Address:    rq.Address,
	}, nil
}

func (h *handler) updateCustomerAddress(ctx context.Context, customerID string, rq updateAddressRq) (customerRs, error) {
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return customerRs{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := gen_sql.New(tx)

	customer, err := q.UpdateCustomerAddress(ctx, gen_sql.UpdateCustomerAddressParams{
		ID:      customerID,
		Address: rq.Address,
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return customerRs{}, err
	}

	err = h.publishEvent(
		uuid.NewString(),
		api.CustomerAddressUpdated,
		api.CustomerAddressUpdatedEvent{
			Address:    rq.Address,
			CustomerID: customerID,
			Timestamp:  time.Now(),
		})
	if err != nil {
		return customerRs{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return customerRs{}, err
	}

	return customerRs{
		CustomerID: customerID,
		Name:       customer.Name,
		Email:      customer.Email,
		Address:    customer.Address,
	}, nil
}

func (h handler) publishEvent(id string, eventType, event any) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		metrics.EventPublishErrors.Inc()
		return err
	}

	msg := amqp.NewMessage(eventData)
	msg.Annotations = map[any]any{
		"event_id":   id,
		"event_type": eventType,
	}

	err = h.producer.Send(amqp.NewMessage(eventData))
	if err != nil {
		metrics.EventPublishErrors.Inc()
		return err
	}

	return nil
}
