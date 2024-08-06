package usecase

import (
	"context"
	"errors"
	db "meight/db/sqlcgen"
	repositoryInterface "meight/repository/interfaces"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Truck struct {
	Database repositoryInterface.Database
}

func NewTruck(database repositoryInterface.Database) *Truck {
	return &Truck{Database: database}
}

func (t *Truck) AddTruck(truck *db.Truck) (db.Truck, error) {

	queries := db.New(t.Database.GetConnectionPool().(*pgxpool.Pool))

	truckParams := db.CreateTruckParams{
		Plate:     truck.Plate,
		MaxWeight: truck.MaxWeight,
	}

	truckDb, err := queries.CreateTruck(context.Background(), truckParams)

	if err != nil {
		log.Error().Msgf("Error creating truck: %v", err)
		return db.Truck{}, errors.New("error creating truck")
	}

	return truckDb, nil
}
