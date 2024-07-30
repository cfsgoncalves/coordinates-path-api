package repository

import (
	"fmt"
	"meight/configuration"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaAccess struct {
	Producer sarama.AsyncProducer
}

func NewKafkaAccess() *KafkaAccess {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_7_0_0

	KAFKA_HOST := configuration.GetEnvAsString("KAFKA_HOST", "localhost")
	KAFKA_PORT := configuration.GetEnvAsInt("KAFKA_PORT", 9094)

	producer, err := sarama.NewAsyncProducer([]string{fmt.Sprintf("%s:%d", KAFKA_HOST, KAFKA_PORT)}, config)
	if err != nil {
		panic(err)
	}

	return &KafkaAccess{Producer: producer}
}

func (k *KafkaAccess) ProduceMessage(topic string, value string) error {

	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}

	k.Producer.Input() <- message

	select {
	case success := <-k.Producer.Successes():
		log.Debug().Msgf("Kafka message was stored on the topic: %s retrieve offset: %s ", topic, success.Topic)
		return nil
	case err := <-k.Producer.Errors():
		log.Error().Msgf("Producing a kafka message yield an error. Error: %s", err)
		return err
	}
}

func (k *KafkaAccess) CloseProducer() error {
	if err := k.Producer.Close(); err != nil {
		log.Error().Msgf("Closing Kafka Producer yield error. Error %s", err)
		return err
	}
	return nil
}
