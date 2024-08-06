package api

import (
	"bytes"
	"context"
	"encoding/json"
	"meight/configuration"
	db "meight/db/sqlcgen"
	repositoryImpl "meight/repository/implementation"
	"meight/usecase"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAssignOrdersToTruck(t *testing.T) {
	godotenv.Load("../../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		ms := repositoryImpl.NewKafkaAccess()
		assert.Nil(t, err)

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "bar",
			Weight:      10,
			Latitude:    2.0,
			Longitude:   3.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "foo",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		jsonArrayId, err := json.Marshal(orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/orderTruck/:truckPlate", bytes.NewBuffer(jsonArrayId))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.AssignOrdersToTruck(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("max_weight_for_the_truck_error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		ms := repositoryImpl.NewKafkaAccess()
		assert.Nil(t, err)

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      30,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		jsonArrayId, err := json.Marshal(orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/orderTruck/:truckPlate", bytes.NewBuffer(jsonArrayId))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.AssignOrdersToTruck(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "orders weight exceed truck max weight")

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("error_truck_plate_does_not_exist", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		ms := repositoryImpl.NewKafkaAccess()
		assert.Nil(t, err)

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		jsonArrayId, err := json.Marshal(orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/v1/orderTruck/:truckPlate", bytes.NewBuffer(jsonArrayId))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: "33"})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.AssignOrdersToTruck(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "inputed truck plate does not exist")

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

func TestUpdateOrderShippingStatus(t *testing.T) {
	godotenv.Load("../../configuration/test.env")

	t.Run("happy_testing", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		// Update order status
		newStatus := "delivered"

		inputBody := struct {
			Status    string `binding:"required" json:"status"`
			OrderCode string `binding:"required" json:"order_code"`
		}{
			Status:    newStatus,
			OrderCode: order2.OrderCode,
		}

		json, err := json.Marshal(inputBody)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, "/v1/orderTruck/:truckPlate", bytes.NewBuffer([]byte(json)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		currentTime := time.Now()
		DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
		date := currentTime.Format(DATE_FORMAT)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Params = append(c.Params, gin.Param{Key: "date", Value: date})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.UpdateOrderShippingStatus(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("invalid_status", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		// Update order status
		newStatus := "fooo"

		inputBody := struct {
			Status    string `binding:"required" json:"status"`
			OrderCode string `binding:"required" json:"order_code"`
		}{
			Status:    newStatus,
			OrderCode: order2.OrderCode,
		}

		json, err := json.Marshal(inputBody)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, "/v1/orderTruck/:truckPlate", bytes.NewBuffer([]byte(json)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		currentTime := time.Now()
		DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
		date := currentTime.Format(DATE_FORMAT)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Params = append(c.Params, gin.Param{Key: "date", Value: date})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.UpdateOrderShippingStatus(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "\"{ Error: invalid status}\"", w.Body.String())

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		// Update order status
		newStatus := "on-route"

		inputBody := struct {
			Status    string `binding:"required" json:"status"`
			OrderCode string `binding:"required" json:"order_code"`
		}{
			Status:    newStatus,
			OrderCode: order2.OrderCode,
		}

		json, err := json.Marshal(inputBody)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, "/v1/orderTruck/:truckPlate", bytes.NewBuffer([]byte(json)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Params = append(c.Params, gin.Param{Key: "date", Value: "12-12-12"})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.UpdateOrderShippingStatus(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "\"{ Error: invalid date format}\"", w.Body.String())

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

func TestGetBestPath(t *testing.T) {
	godotenv.Load("../../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    50.1218,
			Longitude:   8.9298,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    50.1073,
			Longitude:   8.6647,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodGet, "/v1/path/:truckPlate/:date", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})

		currentTime := time.Now()
		DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
		date := currentTime.Format(DATE_FORMAT)

		c.Params = append(c.Params, gin.Param{Key: "date", Value: date})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.GetBestPath(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusOK, w.Code)

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("invalid_data_format", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodGet, "/v1/path/:truckPlate/:date", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Params = append(c.Params, gin.Param{Key: "date", Value: "09-09-24"})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.GetBestPath(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "\"{ 'Error': invalid date format}\"null", w.Body.String())

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG'")
			assert.Nil(t, err)
		})
	})

	t.Run("fail_to_fetch_from_api", func(t *testing.T) {
		os.Setenv("REQUEST_PATH", "http://localhost:8081")
		gin.SetMode(gin.TestMode)

		// Initialize dependencies
		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)

		ms := repositoryImpl.NewKafkaAccess()

		//Add truck
		truck := db.Truck{
			Plate:     "33-66-MG1",
			MaxWeight: 20.5,
		}

		_, err = usecase.NewTruck(newDb).AddTruck(&truck)
		assert.Nil(t, err)

		//Add orders
		order := db.Order{
			OrderCode:   "u123-123-122",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   2.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		orderInsert2 := db.Order{
			OrderCode:   "u123-123-123",
			Weight:      10,
			Latitude:    1.0,
			Longitude:   33.0,
			Description: pgtype.Text{String: "fooo", Valid: true},
		}

		order1, err := usecase.NewOrder(newDb).AddOrder(&order)
		assert.Nil(t, err)
		assert.NotNil(t, order1)

		order2, err := usecase.NewOrder(newDb).AddOrder(&orderInsert2)
		assert.Nil(t, err)
		assert.NotNil(t, order2)

		orderCodeArray := []string{order1.OrderCode, order2.OrderCode}
		assert.Nil(t, err)

		// Assign orders to truck
		err = usecase.NewDistribution(newDb, cache, ms).AssignOrdersToTruck(truck.Plate, orderCodeArray)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodGet, "/v1/path/:truckPlate/:date", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		currentTime := time.Now()
		DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
		date := currentTime.Format(DATE_FORMAT)

		c.Params = append(c.Params, gin.Param{Key: "truckPlate", Value: truck.Plate})
		c.Params = append(c.Params, gin.Param{Key: "date", Value: date})
		c.Request = req

		// Create the API instance
		ordersAPI := NewDistributionApi(cache, newDb, ms)

		// Call the method to be tested
		ordersAPI.GetBestPath(c)

		// Add assertions here to validate the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		t.Cleanup(func() {
			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.order_trucks WHERE truck_plate = '33-66-MG1'")
			assert.Nil(t, err)

			_, err := newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.orders WHERE description = 'fooo'")
			assert.Nil(t, err)

			_, err = newDb.ConnectionPool.Exec(context.Background(), "DELETE FROM meight.public.trucks WHERE plate = '33-66-MG1'")
			assert.Nil(t, err)
		})
	})
}
