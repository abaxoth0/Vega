package postgres

import (
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	FileMetadataTable "vega_file_discovery/packages/infrastrcuture/database/postgres/table/file_metadata"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/transaction"
)

type (
	ConnectionManager   = *connection.Manager
	FileMetadataManager = *FileMetadataTable.Manager
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
		FileMetadataManager: new(FileMetadataTable.Manager),
	}

	executor.Init(connection)
	transaction.Init(connection)

	return driver
}
