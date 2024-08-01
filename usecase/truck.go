package usecase

import (
	"context"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Truck struct {
	Database repositoryInterface.Database
	Cache    repositoryInterface.Cache
}

func NewTruck(cache repositoryInterface.Cache, database repositoryInterface.Database) *Truck {
	return &Truck{Cache: cache, Database: database}
}

func (t *Truck) AddTruck(truck *sqlcgen.Truck) error {

	queries := sqlcgen.New(t.Database.GetConnectionPool().(*pgxpool.Pool))

	foo := sqlcgen.CreateTruckParams{
		Plate:     truck.Plate,
		MaxWeight: truck.MaxWeight,
	}

	_, err := queries.CreateTruck(context.Background(), foo)

	if err != nil {
		log.Error().Msgf("Error creating truck: %v", err)
		return err
	}

	return nil
}
