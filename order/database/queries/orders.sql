-- name: InsertOrder :exec
INSERT INTO orders(id, customer_id, status)
VALUES ($1, $2, $3);

-- name: GetOrderByID :one
SELECT id, customer_id, status FROM orders WHERE id = $1;
