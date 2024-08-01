package api

import (
	"bytes"
	"context"
	"encoding/json"
	repositoryImpl "meight/repository/implementation"
	sqlcgen "meight/sqlc_gen"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAddNewOrder(t *testing.T) {
	godotenv.Load("../../configuration/test.env")
	t.Run("happy_path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		truck := sqlcgen.Order{
			OrderCode:   "123",
			Weight:      3.0,
			Latitude:    0.1,
			Longitude:   0.1,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		truckJson, err := json.Marshal(truck)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/order", bytes.NewBuffer(truckJson))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(cache, newDb)

		// Call the method to be tested
		ordersAPI.AddNewOrder(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)
		})
	})

	t.Run("invalid_json", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		truck := sqlcgen.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 3,
		}

		truckJson, err := json.Marshal(truck)
		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(cache, newDb)

		req, err := http.NewRequest(http.MethodPost, "/v1/order", bytes.NewBuffer(truckJson))
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.AddNewOrder(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, "SQLSTATE 23502", w.Body.String())

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)
		})
	})
}
