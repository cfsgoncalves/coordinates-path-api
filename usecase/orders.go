package usecase

import (
	"context"
	repository "meight/repository/implementation"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
)

type Order struct {
	Database repository.DBAccess
	Cache    repositoryInterface.Cache
}

func NewOrder(cache repositoryInterface.Cache, database repository.DBAccess) *Order {
	return &Order{Cache: cache, Database: database}
}

func (o *Order) AddOrder(order *sqlcgen.Order) error {

	queries := sqlcgen.New(o.Database.ConnectionPool)

	_, err := queries.CreateOrder(context.Background(), sqlcgen.CreateOrderParams{
		Weight:      order.Weight,
		Latitude:    order.Latitude,
		Longitude:   order.Longitude,
		Description: order.Description,
	})

	if err != nil {
		return err
	}

	return nil
}
