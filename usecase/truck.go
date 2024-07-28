package usecase

import (
	"context"
	"meight/entities"
	repository "meight/repository/implementation"
	repositoryInterface "meight/repository/interfaces"
)

type Truck struct {
	Database repository.DBAccess
	Cache    repositoryInterface.Cache
}

func NewTruck(cache repositoryInterface.Cache, database repositoryInterface.Database) *Truck {
	return &Truck{Cache: cache}
}

func (t *Truck) CheckIfTruckPlateExists() bool {
	return true
}

func (t *Truck) AddTruck(truck *entities.Truck) error {
	queries := entities.New(&t.Database.Connection)

	_, err := queries.CreateTruck(context.Background(), entities.CreateTruckParams{
		Plate:     truck.Plate,
		MaxWeight: truck.MaxWeight,
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *Truck) AssignOrdersToTruck(orderTruck []entities.OrderTruck) error {
	return nil
}
