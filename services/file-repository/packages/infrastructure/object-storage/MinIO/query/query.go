package minioquery

import (
	"context"
	"errors"
	"vega/packages/application"
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
	MinIOCommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/minio/minio-go/v7"
)

var Handler FileApplication.QueryHandler = new(defaultQueryHandler)
var storage = MinIOConnection.Manager

type defaultQueryHandler struct {

}

func (m *defaultQueryHandler) GetFileByName(query *FileApplication.GetFileByNameQuery) (*entity.FileStream, error) {
	if !query.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&query.CommandQuery)
	}

	if FileApplication.IsDirectory(query.Path) {
		// TODO need to make archive with all files in directory and send it
		return nil, errors.New("Requested file is directory")
	}
	path, err := FileApplication.FileNameFromPath(query.Path)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(query.Context, query.ContextTimeout)

	if err := MinIOCommon.IsBucketExist(ctx, query.Bucket); err != nil {
		return nil, err
	}

	object, err := storage.Client.GetObject(ctx, query.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}

	return &entity.FileStream{
		Reader: object,
		Size: stat.Size,
		Context: ctx,
		Cancel: cancel,
	}, nil
}

func (m *defaultQueryHandler) SearchFilesByOwner(query *FileApplication.SearchFilesByOwnerQuery) ([]*entity.File, error) {
	return nil, nil
}

