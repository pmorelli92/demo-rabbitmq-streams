package order

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/env"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/features/customer_sync"
)

func Setup(ctx context.Context, e env.Env, db *sqlx.DB, customerChan <-chan customer_sync.CustomerEvent, mux *http.ServeMux) error {
	h := &handler{
		db:           db,
		customerChan: customerChan,
	}

	go h.consumeCustomerEvents(ctx)

	mux.HandleFunc("POST /orders", h.createOrderHandler)

	return nil
}
