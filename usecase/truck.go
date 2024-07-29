package usecase

import (
	"context"
	repository "meight/repository/implementation"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
)

type Truck struct {
	Database repository.DBAccess
	Cache    repositoryInterface.Cache
}

func NewTruck(cache repositoryInterface.Cache, database repository.DBAccess) *Truck {
	return &Truck{Cache: cache, Database: database}
}

func (t *Truck) CheckIfTruckPlateExists() bool {
	return true
}

func (t *Truck) AddTruck(truck *sqlcgen.Truck) error {

	queries := sqlcgen.New(t.Database.ConnectionPool)

	foo := sqlcgen.CreateTruckParams{
		Plate:     truck.Plate,
		MaxWeight: truck.MaxWeight,
	}

	_, err := queries.CreateTruck(context.Background(), foo)

	if err != nil {
		return err
	}

	return nil
}

func (t *Truck) AssignOrdersToTruck(orderTruck []sqlcgen.OrderTruck) error {
	return nil
}
