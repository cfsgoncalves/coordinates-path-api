package main

import (
	"embed"
	"meight/api"
	apiHandlers "meight/api/handlers"
	repositoryImpl "meight/repository/implementation"
	"meight/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

func main() {
	router := gin.Default()

	newDb, err := repositoryImpl.NewDBAccess()

	if err != nil {
		log.Error().Msgf("main(): Error trying to access database")
	}

	err = newDb.MigrateDB(migrationsFS)

	if err != nil {
		log.Error().Msgf("main(): Error trying to execute the migration system")
		return
	}

	cache := repositoryImpl.NewRedis()
	system := usecase.NewSystemMonitoring(cache)
	truckHandlers := apiHandlers.NewTruckApi(cache, newDb)
	distributionHandlers := apiHandlers.NewDistributionApi(cache, newDb)

	api.HTTPRouteEndpoints(router, system, distributionHandlers, truckHandlers)

	router.Run(":8080")

}
