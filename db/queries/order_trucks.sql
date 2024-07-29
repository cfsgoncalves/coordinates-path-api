-- name: ListOrderTrucksByPlateAndDate :many
SELECT date, order_id, truck_plate, latitude, longitude FROM order_trucks, orders
WHERE truck_plate = $1 AND "date" = $2 AND orders.id = order_trucks.order_id
ORDER BY date;

-- name: CreateOrderTrucks :copyfrom
INSERT INTO order_trucks (date, order_id, truck_plate, order_sequence, order_status) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateOrderTrucks :one
UPDATE order_trucks SET order_sequence=$1 WHERE date=$2 AND order_id=$3 AND truck_plate=$4 RETURNING *;

-- name: UpdateOrderTruckStatus :exec
UPDATE order_trucks SET order_status=$1 WHERE date=$2 AND order_id=$3 AND truck_plate=$4;

-- name: ListOrderTrucksByPlateAndDateAndOrderStatus :many
SELECT * FROM order_trucks
WHERE truck_plate = $1 AND "date" = $2 AND order_status = $3
ORDER BY date;




