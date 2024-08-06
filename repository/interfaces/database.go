package repository

import (
	"embed"
)

type Database interface {
	MigrateDB(embed.FS) error
	GetConnectionPool() interface{}
	Ping() bool
}
