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