package usecase

import (
	"fmt"
	"meight/configuration"
	repository "meight/repository/implementation"
	repositoryInterface "meight/repository/interfaces"
	sqlcgen "meight/sqlc_gen"
	"net/http"
	"time"
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

func (d *Distribution) GetBestPath(truckPlate string) ([]sqlcgen.Order, error) {
	return []sqlcgen.Order{}, nil
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
