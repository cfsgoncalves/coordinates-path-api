-- name: ListOrders :many
SELECT * FROM orders
WHERE order_code = $1
ORDER BY order_code;

-- name: CreateOrder :one
INSERT INTO orders (
  order_code, weight, latitude, longitude, description
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOrdersWeightByOrderIds :one
SELECT SUM(orders.weight) as total_weight FROM orders 
WHERE order_code = ANY($1::text[]);

-- name: ListOrdersToBeAssigned :many
SELECT * FROM orders
WHERE NOT exists
(
SELECT order_code FROM order_trucks ot
WHERE order_code  = orders.order_code
and ot.order_status = 'waiting' and Date(ot.date) >= NOW()::timestamp::date
);