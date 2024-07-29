package api

import (
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

func (D *DistributionAPI) GetBestPath(c *gin.Context) {
	value, exists := c.Get("truckPlate")

	if !exists {
		log.Error().Msgf("Bad Request. Variable truckPlate does not exists")
		c.Status(http.StatusBadRequest)
	}

	entitiesOrders, err := D.Distribution.GetBestPath(value.(string))

	if err != nil {
		log.Error().Msgf("GetBestPath yield error. Error: %s", err)
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, entitiesOrders)
}

func (D *DistributionAPI) UpdateShippingOrder(c *gin.Context) {
	request := sqlcgen.OrderTruck{}

	if err := c.BindJSON(&request); err != nil {
		log.Error().Msgf("Bad Request. Could not BindJson to struct")
		c.Status(http.StatusBadRequest)
	}

}
