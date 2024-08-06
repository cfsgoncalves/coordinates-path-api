package repository

import (
	"context"
	"embed"
	"fmt"
	"meight/configuration"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DBAccess struct {
	ConnectionPool   *pgxpool.Pool
	ConnectionString string
}

func NewDBAccess() (*DBAccess, error) {
	ctx := context.Background()

	pgUsername := configuration.GetEnvAsString("DB_USERNAME", "")
	pgPassword := configuration.GetEnvAsString("DB_PASSWORD", "")
	host := configuration.GetEnvAsString("DB_HOST", "")
	port := configuration.GetEnvAsInt("DB_PORT", 5432)
	db := configuration.GetEnvAsString("DB_NAME", "")

	connectionString := fmt.Sprintf("%s:%s@%s:%d/%s?sslmode=disable", pgUsername, pgPassword, host, port, db)

	url := fmt.Sprintf("postgres://" + connectionString)

	connPool, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, err
	}

	err = connPool.Ping(ctx)

	if err != nil {
		log.Error().Msgf("respository.NewDBAccess(): Can't connect to the database")
		return nil, err
	}

	return &DBAccess{ConnectionPool: connPool, ConnectionString: connectionString}, nil
}

func (dB *DBAccess) MigrateDB(migrationsFS embed.FS) error {
	log.Debug().Msg("respository.MigrateDB: Start executing the migration system")

	d, err := iofs.New(migrationsFS, "db/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "postgres://"+dB.ConnectionString)
	if err != nil {
		log.Error().Msgf("respository.MigrateDB():%s", err)
		return err
	}

	err = m.Up()
	if err != nil && err.Error() == "no change" {
		log.Debug().Msg("MigrateDB: There were no new migrations to be executed")
		log.Trace().Msg("MigrateDB: Ended executing migration system")
		return nil
	}

	if err != nil {
		log.Error().Msgf("MigrateDB: Error on the migration. Trying to rollback. Error: %s", err.Error())

		err = m.Down()
		if err != nil {
			log.Error().Msgf("MigrateDB: Error trying to rollback the migration. Error: %s", err.Error())
			return err
		}
		return err
	}

	log.Trace().Msg("MigrateDB: Ended executing migration system")

	return nil
}

func (dB *DBAccess) GetConnectionPool() interface{} {
	return dB.ConnectionPool
}

func (dB *DBAccess) Ping() bool {
	return true
}
