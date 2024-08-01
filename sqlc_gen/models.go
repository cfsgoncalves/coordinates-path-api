// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlcgen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	OrderCode   string      `binding:"required" db:"order_code" json:"order_code"`
	Weight      float64     `binding:"required" db:"weight" json:"weight"`
	Latitude    float64     `binding:"required" db:"latitude" json:"latitude"`
	Longitude   float64     `binding:"required" db:"longitude" json:"longitude"`
	Description pgtype.Text `db:"description" json:"description"`
}

type OrderTruck struct {
	Date          string      `binding:"required" db:"date" json:"date"`
	OrderCode     string      `binding:"required" db:"order_code" json:"order_code"`
	TruckPlate    string      `binding:"required" db:"truck_plate" json:"truck_plate"`
	OrderSequence pgtype.Int4 `binding:"required" db:"order_sequence" json:"order_sequence"`
	OrderStatus   string      `binding:"required" db:"order_status" json:"order_status"`
}

type Truck struct {
	Plate     string  `binding:"required" db:"plate" json:"plate"`
	MaxWeight float64 `binding:"required" db:"max_weight" json:"max_weight"`
}
