package order

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/features/customer_sync"
)

type Order struct {
	ID          string    `db:"id"`
	CustomerID  string    `db:"customer_id"`
	ProductName string    `db:"product_name"`
	Quantity    int       `db:"quantity"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}

func (h *handler) insertOrder(ctx context.Context, tx *sqlx.Tx, order Order) error {
	query := `
		INSERT INTO orders (id, customer_id, product_name, quantity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.ExecContext(ctx, query, order.ID, order.CustomerID, order.ProductName, order.Quantity, order.Price, order.CreatedAt)
	return err
}

func (h *handler) getCustomer(ctx context.Context, customerID string) (customer_sync.Customer, error) {
	var customer customer_sync.Customer
	query := `SELECT id, name, email, address, created_at, updated_at FROM customers WHERE id = $1`

	err := h.db.GetContext(ctx, &customer, query, customerID)
	return customer, err
}
