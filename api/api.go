package api

import (
	apiHandlers "meight/api/handlers"
	usecase "meight/usecase"

	"github.com/gin-gonic/gin"
)

func HTTPRouteEndpoints(router *gin.Engine, system *usecase.SystemMonitoring, distributionApi *apiHandlers.DistributionAPI, truckApi *apiHandlers.TruckAPI) *gin.Engine {
	router.POST("/v1/truck", truckApi.AddNewTruck)
	router.POST("/v1/assignToTruck", truckApi.AssignOrdersToTruck)
	router.GET("/v1/getBestPath/:truckPlate", distributionApi.GetBestPath)
	router.PATCH("/v1/getBestPath/:truckPlate", distributionApi.UpdateShippingOrder)

	router.GET("/health", system.Health)
	router.GET("/ping", system.Ping)

	return router
}
