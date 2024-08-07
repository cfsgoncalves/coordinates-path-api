package api

import (
	"fmt"
	"meight/db/db"
	repository "meight/repository/interfaces"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TruckAPI struct {
	Truck usecase.Truck
}

func NewTruckApi(database repository.Database) *TruckAPI {
	return &TruckAPI{Truck: *usecase.NewTruck(database)}
}

func (t *TruckAPI) AddNewTruck(c *gin.Context) {
	var truck db.Truck

	err := c.BindJSON(&truck)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	truck, err = t.Truck.AddTruck(&truck)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("{ Error: %s}", err))
		return
	}

	c.JSON(http.StatusOK, truck)
}

func (t *TruckAPI) GetTruck(c *gin.Context) {
	plate := c.Param("truckPlate")

	truck, err := t.Truck.GetTruck(plate)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("{ Error: %s}", err))
		return
	}

	c.JSON(http.StatusOK, truck)
}
