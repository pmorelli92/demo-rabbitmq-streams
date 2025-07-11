package customer

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/demo-rabbitmq-streams/api"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/ha"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func Setup(
	db *pgxpool.Pool,
	logger *slog.Logger,
	streamEnv *stream.Environment,
	mux *http.ServeMux) error {

	producer, err := ha.NewReliableProducer(
		streamEnv, api.StreamName,
		stream.NewProducerOptions(),
		func(messageConfirm []*stream.ConfirmationStatus) {})
	if err != nil {
		return err
	}

	httpHandler := httpHandler{
		logger: logger,
		handler: handler{
			db:       db,
			logger:   logger,
			producer: producer,
		},
	}

	mux.HandleFunc("POST /customers", httpHandler.createCustomer)
	mux.HandleFunc("PUT /customers/{id}/address", httpHandler.updateCustomerAddress)

	return nil
}
