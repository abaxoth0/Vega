package objectstorage

import (
	FileApplication "vega_file_repository/packages/application/file"
	minio "vega_file_repository/packages/infrastructure/object-storage/MinIO"
	StorageConnection "vega_file_repository/packages/infrastructure/object-storage/connection"
)

type ObjectStorageDriver interface {
	StorageConnection.Manager
	FileApplication.UseCases
}

var Driver ObjectStorageDriver = minio.InitDriver()
