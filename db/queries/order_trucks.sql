-- name: ListOrderTrucksByPlateAndDate :many
SELECT date, order_trucks.order_code, truck_plate, latitude, longitude, order_status FROM order_trucks, orders
WHERE order_trucks.truck_plate=$1 AND order_trucks.date=$2 AND orders.order_code = order_trucks.order_code
ORDER BY "date";

-- name: CreateOrderTrucks :copyfrom
INSERT INTO order_trucks (date, order_code, truck_plate, order_sequence, order_status) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateOrderTruckSequence :one
UPDATE order_trucks SET order_sequence=$1, order_status=$2 FROM orders 
WHERE date=$3 AND order_trucks.order_code=$4 AND truck_plate=$5 RETURNING date, order_trucks.order_code, truck_plate, order_sequence, order_status, longitude, latitude;

-- name: UpdateOrderTruckStatus :exec
UPDATE order_trucks SET order_status=$1 WHERE date=$2 AND order_code=$3 AND truck_plate=$4;

-- name: GetOrderTrucksByPlateAndDateAndOrderStatus :many
SELECT * FROM order_trucks
WHERE truck_plate = $1 AND "date" = $2 AND order_status = $3
ORDER BY date;

