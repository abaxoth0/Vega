package minio

import (
	FileApplication "vega_file_repository/packages/application/file"
	MinIOCommand "vega_file_repository/packages/infrastructure/object-storage/MinIO/command"
	MinIOConnection "vega_file_repository/packages/infrastructure/object-storage/MinIO/connection"
	MinIOQuery "vega_file_repository/packages/infrastructure/object-storage/MinIO/query"
	StorageConnection "vega_file_repository/packages/infrastructure/object-storage/connection"
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
