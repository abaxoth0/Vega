package minio

import (
	FileApplication "vega/packages/application/file"
	MinIOCommand "vega/packages/infrastructure/object-storage/MinIO/command"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"
	MinIOQuery "vega/packages/infrastructure/object-storage/MinIO/query"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"
)

type Driver struct {
	StorageConnection.Manager
	FileApplication.QueryHandler
	FileApplication.CommandHandler
}

func InitDriver() *Driver {
	return &Driver{
		Manager:        MinIOConnection.Manager,
		QueryHandler:   MinIOQuery.Handler,
		CommandHandler: MinIOCommand.Handler,
	}
}
