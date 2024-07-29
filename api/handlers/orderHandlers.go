package api

import (
	repositoryImpl "meight/repository/implementation"
	repository "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrdersAPI struct {
	Order usecase.Order
}

func NewOrdersApi(cache repository.Cache, database repositoryImpl.DBAccess) *OrdersAPI {
	return &OrdersAPI{Order: *usecase.NewOrder(cache, database)}
}

func (o *OrdersAPI) AddNewOrder(c *gin.Context) {
	var order sqlcgen.Order

	if err := c.BindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := o.Order.AddOrder(&order)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.Status(http.StatusOK)

}
