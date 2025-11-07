package minioquery

import (
	"context"
	"errors"
	"strings"
	"vega/packages/application"
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
	MinIOCommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/abaxoth0/Vega/go-libs/packages/structs"
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
		Size: stat.Size,
		Context: ctx,
		Cancel: cancel,
	}, nil
}

func (h *defaultQueryHandler) GetFileMetadataByPath(query *FileApplication.GetFileByPathQuery) (*entity.FileMetadata, error) {
	if err := h.preprocessQuery(&query.CommandQuery, query.Path); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(query.Context, query.ContextTimeout)
	defer cancel()

	if err := MinIOCommon.IsBucketExist(ctx, query.Bucket); err != nil {
		return nil, err
	}

	stat, err := storage.Client.StatObject(ctx, query.Bucket, query.Path, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	meta := structs.Meta{}

	meta["id"] = stat.UserMetadata["id"]
	meta["original-name"] = stat.UserMetadata["original-name"]
	meta["path"] = stat.Key

	meta["encoding"] = stat.UserMetadata["encoding"]
	meta["mime-type"] = stat.UserMetadata["mime-type"]
	meta["size"] = stat.Size
	meta["checksum"] = stat.ChecksumSHA256
	meta["checksum-type"] = "SHA256"

	meta["owner"] = stat.Owner.ID
	meta["uploaded-by"] = stat.UserMetadata["uploaded-by"]
	meta["permissions"] = stat.UserMetadata["permissions"]

	meta["description"] = stat.UserMetadata["description"]
	meta["categories"] = strings.Split(stat.UserMetadata["categories"], ";")
	meta["tags"] = strings.Split(stat.UserMetadata["tags"], ";")

	meta["status"] = entity.FileStatus(stat.UserMetadata["status"])

	meta["uploaded-at"] = stat.UserMetadata["uploaded-at"]
	meta["updated-at"] = stat.UserMetadata["updated-at"]
	meta["created-at"] = stat.UserMetadata["created-at"]
	meta["accessed-at"] = stat.UserMetadata["accessed-at"]

	metadata, err := entity.NewFileMetadata(meta)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (h *defaultQueryHandler) SearchFilesByOwner(query *FileApplication.SearchFilesByOwnerQuery) ([]*entity.File, error) {
	return nil, nil
}

