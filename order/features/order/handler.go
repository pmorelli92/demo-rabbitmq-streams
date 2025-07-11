package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gen_sql "github.com/pmorelli92/demo-rabbitmq-streams/order/database/generated"
)

type handler struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func (h *handler) createOrder(ctx context.Context, rq createOrderRq) (createOrderRs, error) {
	q := gen_sql.New(h.db)
	customer, err := q.GetCustomerByID(ctx, rq.CustomerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return createOrderRs{}, ErrNotFound
		}
		return createOrderRs{}, err
	}

	// Arbitrary check for testing purposes
	orderStatus := "Processing"
	if customer.Address == "unknown" {
		orderStatus = "Failed"
	}

	orderID := uuid.NewString()
	err = q.InsertOrder(ctx, gen_sql.InsertOrderParams{
		ID:         orderID,
		CustomerID: customer.ID,
		Status:     orderStatus,
	})
	if err != nil {
		return createOrderRs{}, err
	}

	// TODO: We should always create a stream here and publish the events in it

	return createOrderRs{
		OrderID: orderID,
		Status:  orderStatus,
	}, nil
}

func (h *handler) getOrder(ctx context.Context, orderID string) (orderRs, error) {
	q := gen_sql.New(h.db)
	order, err := q.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orderRs{}, ErrNotFound
		}
		return orderRs{}, err
	}

	return orderRs{
		OrderID:    orderID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
	}, nil
}
