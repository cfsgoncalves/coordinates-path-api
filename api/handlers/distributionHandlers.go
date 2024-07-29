package api

import (
	"fmt"
	repositoryImpl "meight/repository/implementation"
	repository "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type DistributionAPI struct {
	Distribution usecase.Distribution
}

func NewDistributionApi(cache repository.Cache, database repositoryImpl.DBAccess) *DistributionAPI {
	return &DistributionAPI{Distribution: *usecase.NewDistribution(database, cache)}
}

func (d *DistributionAPI) GetBestPath(c *gin.Context) {
	value, _ := c.Params.Get("truckPlate")
	date, _ := c.Params.Get("date")

	//Add validation for date format

	entitiesOrders, err := d.Distribution.GetBestPath(value, date)

	if err != nil {
		log.Error().Msgf("GetBestPath yield error. Error: %s", err)
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, entitiesOrders)
}

func (d *DistributionAPI) UpdateShippingOrder(c *gin.Context) {
	value, _ := c.Params.Get("truckPlate")
	date, _ := c.Params.Get("date")

	request := sqlcgen.OrderTruck{}

	if err := c.BindJSON(&request); err != nil {
		log.Error().Msgf("Bad Request. Could not BindJson to struct")
		c.Status(http.StatusBadRequest)
		return
	}

	err := d.Distribution.UpdateShippingOrder(value, date, request.OrderStatus)

	if err != nil {
		log.Error().Msgf("Error UpdateShippingOrder yield error %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ Error: %s}", err))
		return
	}
}

func (d *DistributionAPI) AssignOrdersToTruck(c *gin.Context) {
	value, _ := c.Params.Get("truckPlate")
	request := []sqlcgen.OrderTruck{}

	if err := c.BindJSON(&request); err != nil {
		log.Error().Msgf("Bad Request. Could not BindJson to struct")
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ Error: %s}", "Could not BindJson to struct"))
		return
	}

	err := d.Distribution.AssignOrdersToTruck(value, request)

	if err != nil {
		log.Error().Msgf("Error AssigningOrdersToTruc yield error %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{ Error: %s}", err))
		return
	}

}

func (d *DistributionAPI) GetOrderFromTruck(c *gin.Context) {
	//Get orders from truck where status != delivered
}
