package order

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(
	db *pgxpool.Pool,
	logger *slog.Logger,
	mux *http.ServeMux) {

	httpHandler := httpHandler{
		logger: logger,
		handler: handler{
			db:     db,
			logger: logger,
		},
	}

	mux.HandleFunc("POST /orders", httpHandler.createOrder)
	mux.HandleFunc("GET /orders/{id}", httpHandler.getOrderByID)
}
