package api

import (
	"fmt"
	repository "meight/repository/interfaces"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type DistributionAPI struct {
	Distribution usecase.Distribution
}

func NewDistributionApi(cache repository.Cache, database repository.Database, messageQueue repository.MessageQueue) *DistributionAPI {
	return &DistributionAPI{Distribution: *usecase.NewDistribution(database, cache, messageQueue)}
}

func (d *DistributionAPI) UpdateOrderShippingStatus(c *gin.Context) {
	truckPlate, _ := c.Params.Get("truckPlate")
	date, _ := c.Params.Get("date")

	orderStatus := struct {
		Status    string `binding:"required"`
		OrderCode string `binding:"required"`
	}{}

	if err := c.BindJSON(&orderStatus); err != nil {
		log.Error().Msgf("Bad Request. Could not BindJson to struct")
		c.Status(http.StatusBadRequest)
		return
	}

	err := d.Distribution.UpdateOrderShippingStatus(truckPlate, date, orderStatus.OrderCode, orderStatus.Status)

	if err != nil {
		log.Error().Msgf("Error UpdateShippingOrder yield error %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ Error: %s}", err))
	}
}

func (d *DistributionAPI) AssignOrdersToTruck(c *gin.Context) {
	value, _ := c.Params.Get("truckPlate")
	id := []string{}

	if err := c.BindJSON(&id); err != nil {
		log.Error().Msgf("Bad Request. Could not BindJson to struct")
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ Error: %s}", "Could not BindJson to struct"))
	}

	err := d.Distribution.AssignOrdersToTruck(value, id)

	if err != nil {
		log.Error().Msgf("Error AssigningOrdersToTruc yield error %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ 'Error': %s}", err))
	}
}

func (d *DistributionAPI) GetBestPath(c *gin.Context) {
	value, _ := c.Params.Get("truckPlate")
	date, _ := c.Params.Get("date")

	entitiesOrders, err := d.Distribution.GetBestPath(value, date)

	if err != nil {
		log.Error().Msgf("distributionHandlers.GetBestPath: GetBestPath yield error. Error: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("{ 'Error': %s}", err))
	}

	c.JSON(http.StatusOK, entitiesOrders)
}
