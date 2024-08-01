package api

import (
	repository "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OrdersAPI struct {
	Order usecase.Order
}

func NewOrdersApi(cache repository.Cache, database repository.Database) *OrdersAPI {
	return &OrdersAPI{Order: *usecase.NewOrder(cache, database)}
}

func (o *OrdersAPI) AddNewOrder(c *gin.Context) {
	var order sqlcgen.Order

	err := c.BindJSON(&order)

	if err != nil {
		log.Error().Msgf("Error binding json: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	order, err = o.Order.AddOrder(&order)

	if err != nil {
		log.Error().Msgf("Error adding order: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, order)
}
