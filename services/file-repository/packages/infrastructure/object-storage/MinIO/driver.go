package minio

import (
	FileApplication "vega/packages/application/file"
	MinIOCommandHandler "vega/packages/infrastructure/object-storage/MinIO/command-handler"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"
	MinIOQueryHandler "vega/packages/infrastructure/object-storage/MinIO/query-handler"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"
)

type Driver struct {
	StorageConnection.Manager
	FileApplication.QueryHandler
	FileApplication.CommandHandler
}

func InitDriver() *Driver {
	return &Driver{
		Manager: MinIOConnection.Init(),
		QueryHandler: MinIOQueryHandler.Init(),
		CommandHandler: MinIOCommandHandler.Init(),
	}
}

