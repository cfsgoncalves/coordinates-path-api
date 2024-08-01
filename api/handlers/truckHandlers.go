package api

import (
	repository "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	usecase "meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TruckAPI struct {
	Truck usecase.Truck
}

func NewTruckApi(cache repository.Cache, database repository.Database) *TruckAPI {
	return &TruckAPI{Truck: *usecase.NewTruck(cache, database)}
}

func (t *TruckAPI) AddNewTruck(c *gin.Context) {
	var truck sqlcgen.Truck

	err := c.BindJSON(&truck)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = t.Truck.AddTruck(&truck)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}
