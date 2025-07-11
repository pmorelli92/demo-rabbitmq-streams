-- name: UpsertCustomer :exec
INSERT INTO customers(id, address, updated_at)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET
    address = EXCLUDED.address,
    updated_at = EXCLUDED.updated_at
WHERE EXCLUDED.updated_at > customers.updated_at;
