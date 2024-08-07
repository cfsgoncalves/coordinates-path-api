package api

import (
	apiHandlers "meight/api/handlers"

	"github.com/gin-gonic/gin"
)

func HTTPRouteEndpoints(router *gin.Engine, system *apiHandlers.SystemMonitoringAPI, distributionApi *apiHandlers.DistributionAPI, truckApi *apiHandlers.TruckAPI, orderApi *apiHandlers.OrdersAPI) *gin.Engine {

	router.POST("/v1/truck", truckApi.AddNewTruck)
	router.GET("/v1/truck/:truckPlate", truckApi.GetTruck)
	router.POST("/v1/order", orderApi.AddNewOrder)
	router.GET("/v1/order", orderApi.ListOrdersToBeAssigned)
	router.POST("/v1/orderTruck/:truckPlate", distributionApi.AssignOrdersToTruck)
	router.PATCH("/v1/orderTruck/:truckPlate/:date", distributionApi.UpdateOrderShippingStatus)
	router.GET("/v1/path/:truckPlate/:date", distributionApi.GetBestPath)

	router.GET("/health", system.Health)
	router.GET("/ping", system.Ping)

	return router
}
