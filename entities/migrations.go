package entities

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func NewDBAccess(connectionString string) (*pgx.Conn, error) {
	ctx := context.Background()

	url := fmt.Sprintf("postgres://" + connectionString)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	return conn, nil
}

func MigrateDB(connectionString string, migrationsFS embed.FS) error {
	log.Debug().Msg("MigrateDB: Start executing the migration system")

	d, err := iofs.New(migrationsFS, "db/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "postgres://"+connectionString)
	if err != nil {
		log.Error().Msgf("MigrateDB():%s", err)
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
