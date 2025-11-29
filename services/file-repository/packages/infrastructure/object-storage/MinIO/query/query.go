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

func (h *defaultQueryHandler) preprocessQuery(commandQuery *application.CommandQuery, path string) error {
	if !commandQuery.IsInit() {
		application.InitDefaultCommandQuery(commandQuery)
	}

	err := FileApplication.ValidatePathFormat(path)
	if err != nil {
		return err
	}
	if FileApplication.IsDirectory(path) {
		// TODO need to make archive with all files in directory and send it
		return errors.New("Requested file is directory")
	}

	return nil
}

func (h *defaultQueryHandler) GetFileByPath(query *FileApplication.GetFileByPathQuery) (*entity.FileStream, error) {
	if err := h.preprocessQuery(&query.CommandQuery, query.Path); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(query.Context, query.ContextTimeout)

	if err := MinIOCommon.IsBucketExist(ctx, query.Bucket); err != nil {
		return nil, err
	}

	object, err := storage.Client.GetObject(ctx, query.Bucket, query.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}

	return &entity.FileStream{
		Content: object,
		Size:    stat.Size,
		Context: ctx,
		Cancel:  cancel,
	}, nil
}
