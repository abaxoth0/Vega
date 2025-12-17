package minioquery

import (
	"context"
	"errors"
	FileApplication "vega_file_repository/packages/application/file"
	"vega_file_repository/packages/domain/entity"
	MinIOCommon "vega_file_repository/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega_file_repository/packages/infrastructure/object-storage/MinIO/connection"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
	"github.com/minio/minio-go/v7"
)

var Handler FileApplication.QueryHandler = new(defaultQueryHandler)
var storage = MinIOConnection.Manager

type defaultQueryHandler struct {
}

func (h *defaultQueryHandler) preprocessQuery(commandQuery *cqrs.CommandQuery, path string) error {
	if !commandQuery.IsInit() {
		cqrs.InitDefaultCommandQuery(commandQuery)
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
		if err, ok := err.(minio.ErrorResponse); ok {
			if err.Code == minio.NoSuchKey {
				return nil, errs.StatusNotFound
			}
		}
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		if err, ok := err.(minio.ErrorResponse); ok {
			if err.Code == minio.NoSuchKey {
				return nil, errs.StatusNotFound
			}
		}
		return nil, err
	}

	return &entity.FileStream{
		Content: object,
		Size:    stat.Size,
		Context: ctx,
		Cancel:  cancel,
	}, nil
}
