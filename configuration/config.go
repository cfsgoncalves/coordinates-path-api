package configuration

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func GetEnvAsString(env string, defaultVar string) string {
	envVar := os.Getenv(env)

	if envVar == "" {
		log.Warn().Msgf("Variable %s does not exists. Returning the default variable", env)
		return defaultVar
	}

	return envVar
}

func GetEnvAsInt(env string, defaultVar int) int {
	envVar := os.Getenv(env)

	if envVar == "" {
		log.Warn().Msgf("configuration.GetEnvAsInt():Variable %s does not exists. Returning the default variable", env)
		return defaultVar
	}

	intVar, err := strconv.Atoi(envVar)

	if err != nil {
		log.Error().Msgf("configuration.GetEnvAsInt(): Error converting string to int: %s. Using the default variable", env)
		return defaultVar
	}
	return intVar
}
