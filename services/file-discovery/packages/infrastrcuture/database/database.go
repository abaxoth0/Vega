package DB

import (
	fileapplication "vega_file_discovery/packages/application/file"
	"vega_file_discovery/packages/infrastrcuture/database/postgres"
)


type database interface {
	connector
	fileapplication.UseCases
}

type connector interface {
	Connect() error
	Disconnect() error
}

// Implemets "UseCases" interface of each entity
var Database database = postgres.InitDriver()

type migrate interface {
	Up() error
	Down() error
	Steps(n int) error
}

// Used for applying DB migrations
var Migrate migrate = postgres.Migrate{}
