package usecase

import (
	"context"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Order struct {
	Database repositoryInterface.Database
}

func NewOrder(database repositoryInterface.Database) *Order {
	return &Order{Database: database}
}

func (o *Order) AddOrder(order *sqlcgen.Order) (sqlcgen.Order, error) {

	queries := sqlcgen.New(o.Database.GetConnectionPool().(*pgxpool.Pool))

	orderReturn, err := queries.CreateOrder(context.Background(), sqlcgen.CreateOrderParams{
		OrderCode:   order.OrderCode,
		Weight:      order.Weight,
		Latitude:    order.Latitude,
		Longitude:   order.Longitude,
		Description: order.Description,
	})

	if err != nil {
		log.Error().Msgf("orders.AddOrder: Error adding order: %v", err)
		return sqlcgen.Order{}, err
	}

	return orderReturn, nil
}

func (o *Order) ListOrdersToBeAssigned() ([]sqlcgen.Order, error) {
	queries := sqlcgen.New(o.Database.GetConnectionPool().(*pgxpool.Pool))

	orders, err := queries.ListOrdersToBeAssigned(context.Background())

	if err != nil {
		log.Error().Msgf("orders.AddOrder: Error adding order: %v", err)
		return nil, err
	}
	return orders, nil
}
