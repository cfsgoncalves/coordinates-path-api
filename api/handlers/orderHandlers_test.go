package api

import (
	"bytes"
	"context"
	"encoding/json"
	db "meight/db/sqlcgen"
	repositoryImpl "meight/repository/implementation"
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

		order := db.Order{
			OrderCode:   "123",
			Weight:      3.0,
			Latitude:    0.1,
			Longitude:   0.1,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderJson, err := json.Marshal(order)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/order", bytes.NewBuffer(orderJson))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(newDb)

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

		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 3,
		}

		truckJson, err := json.Marshal(truck)
		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(newDb)

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

func TestListOrdersToBeAssigned(t *testing.T) {
	godotenv.Load("../../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(newDb)
		order, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "foo",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -4.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})

		assert.Nil(t, err)

		order2, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "bar",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -3.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})

		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/v1/order", nil)
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.ListOrdersToBeAssigned(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), order.OrderCode)
		assert.Contains(t, w.Body.String(), order2.OrderCode)

		t.Cleanup(func() {
			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)
		})
	})

	t.Run("get_all_the_order_that_are_not_assigned", func(t *testing.T) {
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(newDb)
		truckAPI := NewTruckApi(newDb)
		order, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "biz",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -4.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})

		assert.Nil(t, err)

		order2, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "bar",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -3.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})
		assert.Nil(t, err)

		_, err = truckAPI.Truck.AddTruck(&db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 56,
		})

		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()
		cache := repositoryImpl.NewRedis()
		d := NewDistributionApi(cache, newDb, ms)

		err = d.Distribution.AssignOrdersToTruck(
			"33-66-MG",
			[]string{order.OrderCode},
		)

		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/v1/order", nil)
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.ListOrdersToBeAssigned(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotContains(t, w.Body.String(), order.OrderCode)
		assert.Contains(t, w.Body.String(), order2.OrderCode)

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("order_waiting_but_from_past_day", func(t *testing.T) {
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		// Create the API instance
		ordersAPI := NewOrdersApi(newDb)
		truckAPI := NewTruckApi(newDb)
		order, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "biz",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -4.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})

		assert.Nil(t, err)

		order2, err := ordersAPI.Order.AddOrder(&db.Order{
			OrderCode:   "bar",
			Weight:      10,
			Latitude:    -2.0,
			Longitude:   -3.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		})
		assert.Nil(t, err)

		_, err = truckAPI.Truck.AddTruck(&db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 56,
		})

		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()
		cache := repositoryImpl.NewRedis()
		d := NewDistributionApi(cache, newDb, ms)

		err = d.Distribution.AssignOrdersToTruck(
			"33-66-MG",
			[]string{order.OrderCode},
		)

		newDb.ConnectionPool.Exec(context.Background(), "UPDATE meight.public.order_trucks SET date = '2021-01-01' WHERE order_code = $1", order.OrderCode)

		assert.Nil(t, err)

		//Create recorder
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/v1/order", nil)
		assert.Nil(t, err)

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the method to be tested
		ordersAPI.ListOrdersToBeAssigned(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), order.OrderCode)
		assert.Contains(t, w.Body.String(), order2.OrderCode)

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})
}
