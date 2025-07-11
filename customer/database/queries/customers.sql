-- name: InsertCustomer :exec
INSERT INTO customers(id, name, email, address, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateCustomerAddress :one
UPDATE customers SET address = $1, updated_at = $2 WHERE id = $3
RETURNING id, name, email, address;
