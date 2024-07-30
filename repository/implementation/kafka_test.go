package repository

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestKafkaInsert(t *testing.T) {

	godotenv.Load("../../configuration/test.env")

	t.Run("happy_path", func(t *testing.T) {
		kafkaAccess := NewKafkaAccess()

		err := kafkaAccess.ProduceMessage("test", "hello")

		assert.Nil(t, err)
	})

}
