-- name: CreateTruck :one
INSERT INTO trucks (
  plate, max_weight
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTruckByPlate :one
SELECT * FROM trucks
WHERE plate = $1;

-- name: GetTruckCapacityByPlateAndDate :one
select plate, max_weight, COALESCE(SUM(weight::float)::float,max_weight::float)::float as capacity from orders, order_trucks, trucks
where orders.order_code = order_trucks.order_code AND trucks.plate = order_trucks.truck_plate
and trucks.plate=$1 and date=$2 and order_trucks.order_status != 'delivered' group by trucks.max_weight, plate;