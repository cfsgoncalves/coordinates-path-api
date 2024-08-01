package repository

type MessageQueue interface {
	ProduceMessage(topic string, value string) error
	CloseProducer() error
}
