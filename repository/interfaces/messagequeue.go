package repository

type MessageQueue interface {
	ProduceMessage(topic string, value string) error
	Ping() bool
}
