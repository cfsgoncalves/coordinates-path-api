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

-- name: ListOrdersByStatus :many
SELECT orders.order_code, weight, latitude, longitude,description FROM orders, order_trucks
WHERE orders.order_code = order_trucks.order_code AND order_trucks.order_status = $1
ORDER BY orders.order_code;