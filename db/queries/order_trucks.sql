-- name: ListOrderTrucksByPlateAndDate :many
SELECT * FROM order_trucks
WHERE truck_plate = $1 AND "date" = $2
ORDER BY date;

-- name: CreateOrderTrucks :copyfrom
INSERT INTO order_trucks (date, order_id, truck_plate, order_sequence) VALUES ($1, $2, $3, $4);



