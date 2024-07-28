// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: orders.sql

package entities

import (
	"context"
)

const listOrders = `-- name: ListOrders :many
SELECT id, weight, latitude, longitude, description FROM orders
WHERE id = $1
ORDER BY id
`

func (q *Queries) ListOrders(ctx context.Context, id int64) ([]Order, error) {
	rows, err := q.db.Query(ctx, listOrders, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.Weight,
			&i.Latitude,
			&i.Longitude,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
