package usecase

import (
	"context"
	"errors"
	"meight/configuration"
	"meight/db/db"
	repositoryInterface "meight/repository/interfaces"
	"time"

	"github.com/jackc/pgx/v5"
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

func (t *Truck) GetTruck(plate string) (db.GetTruckCapacityByPlateAndDateRow, error) {

	queries := db.New(t.Database.GetConnectionPool().(*pgxpool.Pool))

	currentTime := time.Now()
	DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
	date := currentTime.Format(DATE_FORMAT)

	truckDb, err := queries.GetTruckCapacityByPlateAndDate(context.Background(), db.GetTruckCapacityByPlateAndDateParams{
		Plate: plate,
		Date:  date,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		truck, err := queries.GetTruckByPlate(context.Background(), plate)
		if err != nil {
			log.Error().Msgf("usecase.GetTruck: Error getting truck fallback: %v", err)
			return db.GetTruckCapacityByPlateAndDateRow{}, errors.New("error getting truck")
		}

		return db.GetTruckCapacityByPlateAndDateRow{
			Plate:     truck.Plate,
			MaxWeight: truck.MaxWeight,
			Capacity:  truck.MaxWeight,
		}, nil
	}

	if err != nil {
		log.Error().Msgf("usecase.GetTruck: Error getting truck: %v", err)
		return db.GetTruckCapacityByPlateAndDateRow{}, errors.New("error getting truck")
	}

	truckDb.Capacity = truckDb.MaxWeight - truckDb.Capacity

	return truckDb, nil
}
