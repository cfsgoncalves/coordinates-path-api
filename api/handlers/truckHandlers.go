package api

import (
	"meight/entities"
	repository "meight/repository/interfaces"
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
	var truck entities.Truck

	if err := c.BindJSON(&truck); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := t.Truck.AddTruck(&truck)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.Status(http.StatusOK)
}

func (t *TruckAPI) AssignOrdersToTruck(c *gin.Context) {

}
