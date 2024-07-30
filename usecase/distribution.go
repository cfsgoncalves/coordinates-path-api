package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"meight/configuration"
	repository "meight/repository/implementation"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Distribution struct {
	Database repository.DBAccess
	Cache    repositoryInterface.Cache
	HttpC    http.Client
}

func NewDistribution(database repository.DBAccess, cache repositoryInterface.Cache) *Distribution {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	return &Distribution{Database: database, Cache: cache, HttpC: client}
}

func (d *Distribution) GetBestPath(truckPlate string, date string) ([]sqlcgen.OrderTruck, error) {

	queries := sqlcgen.New(d.Database.ConnectionPool)

	orderTruckArray, err := queries.ListOrderTrucksByPlateAndDate(context.Background(), sqlcgen.ListOrderTrucksByPlateAndDateParams{
		TruckPlate: truckPlate,
		Date:       date,
	})

	if err != nil {
		return nil, err
	}

	destinations := []string{}
	for _, value := range orderTruckArray {
		destinations = append(destinations, fmt.Sprintf("%d;%f,%f", value.OrderID, value.Latitude.Float64, value.Longitude.Float64))
	}

	requestPath := getRequestPath(destinations)

	//Do the httpRequest
	resp, err := d.HttpC.Get(requestPath)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	target := HereAPIRequest{}
	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		log.Error().Msgf("Error thing %s", err)
		return nil, err
	}

	//Update on OrderTrucks
	orderTruckArrayToJson := []sqlcgen.OrderTruck{}
	for _, waipoint := range target.Results[0].Waipoints {
		orderTruck, err := queries.UpdateOrderTrucks(context.Background(), sqlcgen.UpdateOrderTrucksParams{
			OrderSequence: waipoint.Id,
			Date:          date,
			OrderID:       int64(waipoint.Sequence),
			TruckPlate:    truckPlate,
		})

		if err != nil {
			return nil, err
		}

		orderTruckArrayToJson = append(orderTruckArrayToJson, orderTruck)
	}

	return orderTruckArrayToJson, nil

}

// Aux method to get the best path
func getRequestPath(destinations []string) string {
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

func (d *Distribution) AssignOrdersToTruck(truckPlate string, orders []sqlcgen.OrderTruck) error {

	orderTruckArray := []sqlcgen.CreateOrderTrucksParams{}
	queries := sqlcgen.New(d.Database.ConnectionPool)

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	for _, value := range orders {
		orderTruckParam := sqlcgen.CreateOrderTrucksParams{
			Date:          date,
			OrderID:       value.OrderID,
			TruckPlate:    truckPlate,
			OrderSequence: pgtype.Int4{},
			OrderStatus:   "WAITING",
		}
		orderTruckArray = append(orderTruckArray, orderTruckParam)
	}

	_, err := queries.CreateOrderTrucks(context.Background(), orderTruckArray)

	if err != nil {
		return err
	}

	return nil
}

func (d *Distribution) UpdateShippingOrder(truckPlate string, date string, status string) error {
	//Update Order Truck on the database

	//Send an event to kafka for each update

	return nil
}
