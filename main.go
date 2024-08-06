package main

import (
	"embed"
	"fmt"
	"meight/api"
	apiHandlers "meight/api/handlers"
	"meight/configuration"
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
		return
	}

	err = newDb.MigrateDB(migrationsFS)

	if err != nil {
		log.Error().Msgf("main(): Error trying to execute the migration system")
		return
	}

	//Initialize dependencies
	cache := repositoryImpl.NewRedis()
	messageQueue := repositoryImpl.NewKafkaAccess()

	system := usecase.NewSystemMonitoring(cache, newDb, messageQueue)

	truckHandlers := apiHandlers.NewTruckApi(newDb)
	distributionHandlers := apiHandlers.NewDistributionApi(cache, newDb, messageQueue)
	orderHandlers := apiHandlers.NewOrdersApi(newDb)

	api.HTTPRouteEndpoints(router, system, distributionHandlers, truckHandlers, orderHandlers)

	SERVER_PORT := configuration.GetEnvAsString("SERVER_PORT", "8081")
	router.Run(fmt.Sprintf(":%s", SERVER_PORT))

}
