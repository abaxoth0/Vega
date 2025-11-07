package objectstorage

import (
	FileApplication "vega/packages/application/file"
	minio "vega/packages/infrastructure/object-storage/MinIO"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"
)

type ObjectStorageDriver interface {
	StorageConnection.Manager
	FileApplication.UseCases
}

var Driver ObjectStorageDriver = minio.InitDriver()
