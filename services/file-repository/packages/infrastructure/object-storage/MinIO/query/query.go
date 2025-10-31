package minioquery

import (
	"context"
	"errors"
	"time"
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
	miniocommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/minio/minio-go/v7"
)

var Handler FileApplication.QueryHandler = new(defaultQueryHandler)
var storage = MinIOConnection.Manager

type defaultQueryHandler struct {

}

func (m *defaultQueryHandler) GetFileByName(query *FileApplication.GetFileByNameQuery) (*entity.FileStream, error) {
	if query.Context == nil {
		query.Context = context.Background()
	}
	if query.ContextTimeout <= 0 {
		// TODO move right part into the constant
		query.ContextTimeout = time.Second * 3
	}

	ctx, cancel := context.WithTimeout(query.Context, query.ContextTimeout)

	ok, err := storage.Client.BucketExists(ctx, query.Bucket)
	if err != nil {
		// TODO temp
		return nil, err
	}
	if !ok {
		// TODO temp error (make it sentinel error for proper handling)
		return nil, errors.New("Bucket doesn't exist")
	}

	if miniocommon.IsDirectory(query.Path) {
		// TODO need to make archive with all files in directory and send it
		return nil, errors.New("Requested file is directory")
	}

	object, err := storage.Client.GetObject(ctx, query.Bucket, query.FileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &entity.FileStream{
		Reader: object,
		Context: ctx,
		Cancel: cancel,
	}, nil
}

func (m *defaultQueryHandler) SearchFilesByOwner(query *FileApplication.SearchFilesByOwnerQuery) ([]*entity.File, error) {
	return nil, nil
}

