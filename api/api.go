package api

import (
	apiHandlers "meight/api/handlers"
	usecase "meight/usecase"

	"github.com/gin-gonic/gin"
)

func HTTPRouteEndpoints(router *gin.Engine, system *usecase.SystemMonitoring, distributionApi *apiHandlers.DistributionAPI, truckApi *apiHandlers.TruckAPI, orderApi *apiHandlers.OrdersAPI) *gin.Engine {
	router.POST("/v1/truck", truckApi.AddNewTruck)
	router.POST("/v1/order", orderApi.AddNewOrder)
	router.POST("/v1/assignToTruck/:truckPlate", distributionApi.AssignOrdersToTruck)
	router.GET("/v1/orderFromTruck/:truckPlate", distributionApi.GetOrderFromTruck)
	router.GET("/v1/getBestPath/:truckPlate/:date", distributionApi.GetBestPath)
	router.PATCH("/v1/getBestPath/:truckPlate/:date", distributionApi.UpdateShippingOrder)

	router.GET("/health", system.Health)
	router.GET("/ping", system.Ping)

	return router
}
