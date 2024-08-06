package usecase

import (
	"context"
	db "meight/db/sqlcgen"
	repositoryInterface "meight/repository/interfaces"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Order struct {
	Database repositoryInterface.Database
}

func NewOrder(database repositoryInterface.Database) *Order {
	return &Order{Database: database}
}

func (o *Order) AddOrder(order *db.Order) (db.Order, error) {

	queries := db.New(o.Database.GetConnectionPool().(*pgxpool.Pool))

	orderReturn, err := queries.CreateOrder(context.Background(), db.CreateOrderParams{
		OrderCode:   order.OrderCode,
		Weight:      order.Weight,
		Latitude:    order.Latitude,
		Longitude:   order.Longitude,
		Description: order.Description,
	})

	if err != nil {
		log.Error().Msgf("orders.AddOrder: Error adding order: %v", err)
		return db.Order{}, err
	}

	return orderReturn, nil
}

func (o *Order) ListOrdersToBeAssigned() ([]db.Order, error) {
	queries := db.New(o.Database.GetConnectionPool().(*pgxpool.Pool))

	orders, err := queries.ListOrdersToBeAssigned(context.Background())

	if err != nil {
		log.Error().Msgf("orders.AddOrder: Error adding order: %v", err)
		return nil, err
	}
	return orders, nil
}
