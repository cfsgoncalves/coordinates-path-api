package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"meight/configuration"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"net/http"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var ALLOWED_STATUS = []string{"waiting", "on-going", "delivered"}

type Distribution struct {
	Database     repositoryInterface.Database
	Cache        repositoryInterface.Cache
	HttpC        http.Client
	MessageQueue repositoryInterface.MessageQueue
}

func NewDistribution(database repositoryInterface.Database, cache repositoryInterface.Cache, messageQueue repositoryInterface.MessageQueue) *Distribution {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return &Distribution{Database: database, Cache: cache, HttpC: client, MessageQueue: messageQueue}
}

func (d *Distribution) GetBestPath(truckPlate string, date string) ([]sqlcgen.UpdateOrderTruckSequenceRow, error) {
	// Validation for date format
	DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
	_, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	//Check if data exists on cache
	orderTruckArrayToJson := d.isDataOnCache(truckPlate, date)
	if d.isDataOnCache(truckPlate, date) != nil {
		log.Debug().Msgf("usecase.GetBestPath: Data found on cache")
		return orderTruckArrayToJson, nil
	}

	queries := sqlcgen.New(d.Database.GetConnectionPool().(*pgxpool.Pool))

	orderTruckArray, err := queries.ListOrderTrucksByPlateAndDate(context.Background(), sqlcgen.ListOrderTrucksByPlateAndDateParams{
		TruckPlate: truckPlate,
		Date:       date,
	})

	if err != nil {
		log.Error().Msgf("usecase.GetBestPath: Error getting order trucks by plate and date: %v", err)
		return nil, err
	}

	destinations := []string{}
	for _, value := range orderTruckArray {
		destinations = append(destinations, fmt.Sprintf("%s;%f,%f", value.OrderCode, value.Latitude, value.Longitude))
	}

	//Build the request path
	requestPath := d.getRequestPath(destinations)

	//Check for coordinate points

	//Do the httpRequest
	target, err := d.fetchFromOrderingAPI(requestPath)

	if err != nil {
		log.Error().Msgf("usecase.GetBestPath: Error decoding response body %s", err)
		return nil, err
	}

	//Update on OrderTrucks OrderSequence and OrderStatus
	for _, waipoint := range target.Results[0].Waipoints[1:] {

		orderTruck, err := queries.UpdateOrderTruckSequence(context.Background(), sqlcgen.UpdateOrderTruckSequenceParams{
			OrderSequence: waipoint.Sequence,
			OrderStatus:   "on-going",
			Date:          date,
			OrderCode:     waipoint.Id,
			TruckPlate:    truckPlate,
		})

		if err != nil {
			log.Error().Msgf("usecase.GetBestPath: Error updating order truck sequence: %v", err)
			return nil, err
		}

		orderTruckArrayToJson = append(orderTruckArrayToJson, orderTruck)
	}

	//Save to cache
	err = d.saveToCache(truckPlate, date, orderTruckArrayToJson)

	if err != nil {
		log.Error().Msgf("usecase.GetBestPath: Error saving to cache: %v", err)
	}

	return orderTruckArrayToJson, nil
}

func (d *Distribution) AssignOrdersToTruck(truckPlate string, orderCodes []string) error {
	orderTruckArray := []sqlcgen.CreateOrderTrucksParams{}
	queries := sqlcgen.New(d.Database.GetConnectionPool().(*pgxpool.Pool))

	ordersWeight, err := d.getOrderWeightBasedOnIds(orderCodes)

	if err != nil {
		log.Error().Msgf("distribution.AssignOrdersToTruck: Error GetOrderWeightBasedOnIds yield error %s", err)
		return err
	}

	truck, err := queries.GetTruckByPlate(context.Background(), truckPlate)

	if err != nil {
		log.Error().Msgf("distribution.AssignOrdersToTruck: Error GetTruckByPlate yield error %s", err)
		return errors.New("inputed truck plate does not exist")
	}

	//weight validation
	if ordersWeight > int64(truck.MaxWeight) {
		log.Error().Msgf("distribution.AssignOrdersToTruck: Error AssigningOrdersToTruck. Yield error orders weight exceed truck max weight")
		return errors.New("yield error orders weight exceed truck max weight")
	}

	// Generate date for the assignment
	currentTime := time.Now()
	DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")
	date := currentTime.Format(DATE_FORMAT)

	for _, orderCode := range orderCodes {
		orderTruckParam := sqlcgen.CreateOrderTrucksParams{
			Date:          date,
			OrderCode:     orderCode,
			TruckPlate:    truckPlate,
			OrderSequence: pgtype.Int4{},
			OrderStatus:   "waiting",
		}
		orderTruckArray = append(orderTruckArray, orderTruckParam)
	}

	_, err = queries.CreateOrderTrucks(context.Background(), orderTruckArray)

	if err != nil {
		log.Error().Msgf("distribution.AssignOrdersToTruck: Error creating orderTrucks: %v", err)
		return err
	}

	//Send an event to kafka for each creation
	for _, orderTruck := range orderTruckArray {
		err = d.sendMessageToQueue(orderTruck.OrderCode, orderTruck.OrderStatus)
		if err != nil {
			log.Error().Msgf("Error sending message to message queue with orderCode: %s. Yield error: %v", orderTruck.OrderCode, err)
		}
	}

	return nil
}

func (d *Distribution) UpdateOrderShippingStatus(truckPlate string, date string, orderCode string, status string) error {
	// Validation for status
	if !slices.Contains(ALLOWED_STATUS, status) {
		return errors.New("invalid status")
	}

	// Validation for date format
	DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")

	//Validate date paramenter
	_, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return errors.New("invalid date format")
	}

	queries := sqlcgen.New(d.Database.GetConnectionPool().(*pgxpool.Pool))

	err = queries.UpdateOrderTruckStatus(context.Background(), sqlcgen.UpdateOrderTruckStatusParams{
		OrderStatus: status,
		Date:        date,
		OrderCode:   orderCode,
		TruckPlate:  truckPlate,
	})

	if err != nil {
		log.Error().Msgf("Error updating order truck status: %v", err)
		return err
	}

	// Send Message to Queue
	err = d.sendMessageToQueue(orderCode, status)

	if err != nil {
		log.Error().Msgf("Error sending message to message queue with orderCode: %s. Yield error: %v", orderCode, err)
		return errors.New("error sending message to message queue")
	}

	// Get from cache
	orderTruckArrayToJson := d.isDataOnCache(truckPlate, date)

	for index, value := range orderTruckArrayToJson {
		if value.OrderCode == orderCode {
			orderTruckArrayToJson[index].OrderStatus = status
		}
	}

	err = d.saveToCache(truckPlate, date, orderTruckArrayToJson)

	if err != nil {
		log.Error().Msgf("Error saving to cache: %v", err)
	}

	return nil
}

func (d *Distribution) GetOrdersFromTruck(truckPlate string, date string, status string) ([]sqlcgen.ListOrderTrucksByPlateAndDateRow, error) {
	// Validation for status
	if !slices.Contains(ALLOWED_STATUS, status) {
		return nil, errors.New("invalid status")
	}

	// Validation for date format
	DATE_FORMAT := configuration.GetEnvAsString("DATE_FORMAT", "2006-01-02")

	//Validate date paramenter
	_, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	queries := sqlcgen.New(d.Database.GetConnectionPool().(*pgxpool.Pool))

	orderTruckArray, err := queries.ListOrderTrucksByPlateAndDate(context.Background(), sqlcgen.ListOrderTrucksByPlateAndDateParams{
		TruckPlate: truckPlate,
		Date:       date,
	})

	if err != nil {
		log.Error().Msgf("Error getting order trucks by plate and date: %v", err)
		return nil, errors.New("error trying to acess the database")
	}

	orderTruckArrayStatus := []sqlcgen.ListOrderTrucksByPlateAndDateRow{}
	for _, value := range orderTruckArray {
		if value.OrderStatus == status {
			orderTruckArrayStatus = append(orderTruckArrayStatus, value)
		}
	}

	return orderTruckArrayStatus, nil
}

func (d *Distribution) getOrderWeightBasedOnIds(orderIds []string) (int64, error) {
	queries := sqlcgen.New(d.Database.GetConnectionPool().(*pgxpool.Pool))

	orders, err := queries.GetOrdersWeightByOrderIds(context.Background(), orderIds)
	if err != nil {
		log.Error().Msgf("Error getting orders by ids: %v", err)
		return -1, err
	}
	return orders, nil
}

func (d *Distribution) getRequestPath(destinations []string) string {
	if len(destinations) < 1 {
		return ""
	}

	REQUEST_PATH := configuration.GetEnvAsString("REQUEST_PATH", "https://wps.hereapi.com/v8/findsequence2")
	API_KEY := configuration.GetEnvAsString("API_KEY", "Xi-sfj72ReKKh6O1_r0oTz9AUVPx5j84JzMaeRj-mb8")
	STARTING_POINT := configuration.GetEnvAsString("STARTING_POINT", "FintechHouse;38.71814,-9.14552")
	CALCULATION_MODE := configuration.GetEnvAsString("CALCULATION_MODE", "shortest;truck;traffic:disabled")
	request_final := fmt.Sprintf("%s?apiKey=%s&start=%s&mode=%s", REQUEST_PATH, API_KEY, STARTING_POINT, CALCULATION_MODE)

	destinationString := ""
	for index, value := range destinations {
		destinationString += fmt.Sprintf("&destination%d=%s", index, value)
	}

	request_final = request_final + destinationString

	return request_final
}

func (d *Distribution) isDataOnCache(truckPlate string, date string) []sqlcgen.UpdateOrderTruckSequenceRow {
	//Check if data exists on cache
	orderTruckArrayCache, err := d.Cache.Get(fmt.Sprintf("%s-%s", truckPlate, date))
	if err != nil {
		log.Error().Msgf("usecase.GetBestPath(): Error yield acessing cache service. Error: %s", err)
	}

	// Return the values on cache with 200 and value
	var orderTruck []sqlcgen.UpdateOrderTruckSequenceRow
	if orderTruckArrayCache != "" {
		err := json.Unmarshal([]byte(orderTruckArrayCache), &orderTruck)

		if err != nil {
			log.Error().Msgf("usecase.GetBestPath(): Error yield while unmarshaling cache value. Error: %s", err)
			return nil
		}
		return orderTruck
	}

	return nil
}

func (d *Distribution) saveToCache(truckPlate string, date string, orderTruckArrayToJson []sqlcgen.UpdateOrderTruckSequenceRow) error {
	jsonString, err := json.Marshal(orderTruckArrayToJson)

	if err != nil {
		log.Error().Msgf("usecase.saveToCache: Error marshalling orderTruckArrayToJson: %v", err)
		return err
	}

	err = d.Cache.Insert(fmt.Sprintf("%s-%s", truckPlate, date), string(jsonString))

	if err != nil {
		log.Error().Msgf("usecase.saveToCache: Error while inserting to cache: %v", err)
		return err
	}

	return nil
}

func (d *Distribution) fetchFromOrderingAPI(requestPath string) (HereAPIRequest, error) {
	resp, err := d.HttpC.Get(requestPath)

	if err != nil {
		log.Error().Msgf("usecase.fetchFromOrderingAPI: Error getting http request %s", err)
		return HereAPIRequest{}, err
	}

	defer resp.Body.Close()

	target := HereAPIRequest{}
	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		log.Error().Msgf("usecase.fetchFromOrderingAPI: Error decoding response body %s", err)
		return HereAPIRequest{}, err
	}

	return target, nil
}

func (d *Distribution) sendMessageToQueue(orderCode string, status string) error {
	message := struct {
		OrderCode string
		Status    string
		Message   string
	}{OrderCode: orderCode, Status: status, Message: fmt.Sprintf("Order %s status updated to %s", orderCode, status)}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Debug().Msgf("Error marshalling message: %v", err)
		return err
	}

	//Send an event to kafka for each update
	err = d.MessageQueue.ProduceMessage("order_truck", string(jsonMessage))

	if err != nil {
		log.Error().Msgf("Error producing message: %v", err)
		return err
	}

	return nil
}
