package customer_sync

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomerEvent struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	CustomerID string      `json:"customer_id"`
	Timestamp  time.Time   `json:"timestamp"`
	Data       interface{} `json:"data"`
}

type Customer struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Address   string    `db:"address"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type handler struct {
	db *sqlx.DB
}

func (h *handler) upsertCustomer(ctx context.Context, customer Customer) error {
	query := `
		INSERT INTO customers (id, name, email, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			address = EXCLUDED.address,
			updated_at = EXCLUDED.updated_at
	`
	
	_, err := h.db.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email, customer.Address, customer.CreatedAt, customer.UpdatedAt)
	return err
}

func (h *handler) updateCustomerAddress(ctx context.Context, customerID, address string) error {
	query := `UPDATE customers SET address = $1, updated_at = $2 WHERE id = $3`
	
	_, err := h.db.ExecContext(ctx, query, address, time.Now(), customerID)
	return err
}