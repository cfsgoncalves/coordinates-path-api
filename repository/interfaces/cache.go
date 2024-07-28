package repository

type Cache interface {
	Insert(string, string) error
	Get(string) (string, error)
	Ping() bool
}
