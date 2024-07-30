package repository

type MessageQueue interface {
	ProduceMessage(topic string, value string)
	CloseProducer()
}
