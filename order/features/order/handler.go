package order

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/features/customer_sync"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/metrics"
)

type handler struct {
	db           *sqlx.DB
	customerChan <-chan customer_sync.CustomerEvent
}

type CreateOrderRequest struct {
	CustomerID  string  `json:"customer_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

func (h *handler) consumeCustomerEvents(ctx context.Context) {
	for {
		select {
		case event := <-h.customerChan:
			log.Printf("Received customer event: %s for customer %s", event.EventType, event.CustomerID)
		case <-ctx.Done():
			log.Println("Stopping customer event consumption")
			return
		}
	}
}

func (h *handler) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues("POST", "/orders").Observe(time.Since(start).Seconds())
	}()

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.HTTPRequests.WithLabelValues("POST", "/orders", "400").Inc()
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	customer, err := h.getCustomer(ctx, req.CustomerID)
	if err != nil {
		metrics.HTTPRequests.WithLabelValues("POST", "/orders", "404").Inc()
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	orderID := uuid.New().String()
	order := Order{
		ID:          orderID,
		CustomerID:  req.CustomerID,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CreatedAt:   time.Now(),
	}

	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		metrics.HTTPRequests.WithLabelValues("POST", "/orders", "500").Inc()
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func() { _ = tx.Rollback() }()

	err = h.insertOrder(ctx, tx, order)
	if err != nil {
		metrics.HTTPRequests.WithLabelValues("POST", "/orders", "500").Inc()
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		metrics.HTTPRequests.WithLabelValues("POST", "/orders", "500").Inc()
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	metrics.OrdersCreated.Inc()
	metrics.HTTPRequests.WithLabelValues("POST", "/orders", "201").Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id":      orderID,
		"customer_id":   req.CustomerID,
		"customer_name": customer.Name,
		"product_name":  req.ProductName,
		"quantity":      req.Quantity,
		"price":         req.Price,
		"status":        "created",
	})
}
