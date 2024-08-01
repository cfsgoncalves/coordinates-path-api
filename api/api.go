package api

import (
	apiHandlers "meight/api/handlers"
	usecase "meight/usecase"

	"github.com/gin-gonic/gin"
)

func HTTPRouteEndpoints(router *gin.Engine, system *usecase.SystemMonitoring, distributionApi *apiHandlers.DistributionAPI, truckApi *apiHandlers.TruckAPI, orderApi *apiHandlers.OrdersAPI) *gin.Engine {

	router.POST("/v1/truck", truckApi.AddNewTruck)
	router.POST("/v1/order", orderApi.AddNewOrder)
	router.GET("/v1/order/:status", distributionApi.GetTruckOrders)
	router.POST("/v1/orderTruck/:truckPlate", distributionApi.AssignOrdersToTruck)
	router.PATCH("/v1/orderTruck/:truckPlate/:date", distributionApi.UpdateOrderShippingStatus)
	router.GET("/v1/getBestPath/:truckPlate/:date", distributionApi.GetBestPath)

	router.GET("/health", system.Health)
	router.GET("/ping", system.Ping)

	return router
}
