package repository

import (
	"fmt"
	"meight/configuration"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaAccess struct {
	Config *sarama.Config
}

func NewKafkaAccess() *KafkaAccess {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_7_0_0

	return &KafkaAccess{Config: config}
}

func (k *KafkaAccess) ProduceMessage(topic string, value string) error {

	KAFKA_HOST := configuration.GetEnvAsString("KAFKA_HOST", "localhost")
	KAFKA_PORT := configuration.GetEnvAsInt("KAFKA_PORT", 9094)

	producer, err := sarama.NewAsyncProducer([]string{fmt.Sprintf("%s:%d", KAFKA_HOST, KAFKA_PORT)}, k.Config)

	defer producer.Close()

	if err != nil {
		log.Error().Msgf("Error creating Kafka producer: %v", err)
		return nil
	}

	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}

	producer.Input() <- message

	select {
	case success := <-producer.Successes():
		log.Debug().Msgf("Kafka message was stored on the topic: %s retrieve offset: %s ", topic, success.Topic)
		return nil
	case err := <-producer.Errors():
		log.Error().Msgf("Producing a kafka message yield an error. Error: %s", err)
		return err
	}
}

func (k *KafkaAccess) Ping() bool {
	return true
}
