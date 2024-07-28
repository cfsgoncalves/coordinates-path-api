-- name: ListOrders :many
SELECT * FROM orders
WHERE id = $1
ORDER BY id;