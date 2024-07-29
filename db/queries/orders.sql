-- name: ListOrders :many
SELECT * FROM orders
WHERE id = $1
ORDER BY id;

-- name: CreateOrder :one
INSERT INTO orders (
  weight, latitude, longitude, description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

