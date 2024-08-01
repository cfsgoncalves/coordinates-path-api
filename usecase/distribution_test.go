package usecase

import (
	"os"
	"testing"

	repositoryImpl "meight/repository/implementation"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestPath(t *testing.T) {
	godotenv.Load("../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {

		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)
		ms := repositoryImpl.NewKafkaAccess()

		d := NewDistribution(newDb, cache, ms)

		os.Setenv("REQUEST_PATH", "https://wps.hereapi.com/v8/findsequence2")
		os.Setenv("API_KEY", "Xi-sfj72ReKKh6O1_r0oTz9AUVPx5j84JzMaeRj-mb8")
		os.Setenv("STARTING_POINT", "FintechHouse;38.71814,-9.14552")
		os.Setenv("CALCULATION_MODE", "shortest;truck;traffic:disabled")
		destinations := []string{"50.1073,8.6647", "49.8728,8.6326", "50.0505,8.5698", "50.1218,8.9298"}

		value := d.getRequestPath(destinations)

		assert.Contains(t, value, "https://wps.hereapi.com/v8/findsequence2?apiKey=Xi-sfj72ReKKh6O1_r0oTz9AUVPx5j84JzMaeRj-mb8&start=FintechHouse;38.71814,-9.14552&mode=shortest;truck;traffic:disabled")
		assert.Contains(t, value, "&destination0=50.1073,8.6647&destination1=49.8728,8.6326&destination2=50.0505,8.5698&destination3=50.1218,8.9298")
	})

	t.Run("env_variables_not_defined", func(t *testing.T) {
		os.Unsetenv("REQUEST_PATH")
		os.Unsetenv("API_KEY")
		os.Unsetenv("STARTING_POINT")
		os.Unsetenv("CALCULATION_MODE")

		cache := repositoryImpl.NewRedis()
		newDb, err := repositoryImpl.NewDBAccess()
		assert.Nil(t, err)
		ms := repositoryImpl.NewKafkaAccess()

		d := NewDistribution(newDb, cache, ms)

		destinations := []string{"50.1073,8.6647", "49.8728,8.6326", "50.0505,8.5698", "50.1218,8.9298"}

		value := d.getRequestPath(destinations)

		assert.Contains(t, value, "https://wps.hereapi.com/v8/findsequence2?apiKey=Xi-sfj72ReKKh6O1_r0oTz9AUVPx5j84JzMaeRj-mb8&start=FintechHouse;38.71814,-9.14552&mode=shortest;truck;traffic:disabled")
		assert.Contains(t, value, "&destination0=50.1073,8.6647&destination1=49.8728,8.6326&destination2=50.0505,8.5698&destination3=50.1218,8.9298")

	})
}
