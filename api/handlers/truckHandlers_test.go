package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"meight/db/db"
	repositoryImpl "meight/repository/implementation"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewTruckApi(t *testing.T) {
	godotenv.Load("../../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 3,
		}

		truckJson, err := json.Marshal(truck)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/truck", bytes.NewBuffer(truckJson))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewTruckApi(newDb)

		// Call the method to be tested
		ordersAPI.AddNewTruck(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "33-66-MG")

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("invalid_json", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		order := db.Order{
			Weight:      2,
			Latitude:    0,
			Longitude:   0,
			Description: pgtype.Text{},
		}

		truckJson, err := json.Marshal(order)
		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewTruckApi(newDb)

		req, err := http.NewRequest(http.MethodPost, "/v1/truck", bytes.NewBuffer(truckJson))
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.AddNewTruck(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, "SQLSTATE 23502", w.Body.String())

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("duplicate_key_error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 3,
		}

		db, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		_, err = NewTruckApi(db).Truck.AddTruck(&truck)
		assert.Nil(t, err)

		truckJson, err := json.Marshal(truck)
		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewTruckApi(newDb)

		req, err := http.NewRequest(http.MethodPost, "/v1/truck", bytes.NewBuffer(truckJson))
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.AddNewTruck(c)

		// Add assertions here to validate the response
		fmt.Println(w.Body.String())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error creating truck")

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})
}
