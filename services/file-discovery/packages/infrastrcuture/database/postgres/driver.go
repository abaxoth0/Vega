package postgres

import (
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	filemetatable "vega_file_discovery/packages/infrastrcuture/database/postgres/table/file_meta"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/transaction"
)

type (
	ConnectionManager   = *connection.Manager
	FileMetadataManager = *filemetatable.Manager
)

type postgers struct {
	ConnectionManager
	FileMetadataManager
}

var driver *postgers

func InitDriver() *postgers {
	connection := new(connection.Manager)

	driver = &postgers{
		ConnectionManager: connection,
		FileMetadataManager: new(filemetatable.Manager),
	}

	executor.Init(connection)
	transaction.Init(connection)

	return driver
}
