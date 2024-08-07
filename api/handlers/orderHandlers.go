package api

import (
	"meight/db/db"
	repository "meight/repository/interfaces"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OrdersAPI struct {
	Order usecase.Order
}

func NewOrdersApi(database repository.Database) *OrdersAPI {
	return &OrdersAPI{Order: *usecase.NewOrder(database)}
}

func (o *OrdersAPI) AddNewOrder(c *gin.Context) {
	var order db.Order

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

func (o *OrdersAPI) ListOrdersToBeAssigned(c *gin.Context) {
	orders, err := o.Order.ListOrdersToBeAssigned()
	if err != nil {
		log.Error().Msgf("Error getting orders: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, orders)
}
