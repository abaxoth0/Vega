package minioquery

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"vega/packages/application"
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
	MinIOCommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/abaxoth0/Vega/libs/go/packages/structs"
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

func splitGroupedField(s string) []string {
	// If string is empty strings.Split() retuns slice with 1 element - empty string.
	split := strings.Split(s, ";")
	if len(split) == 1 && split[0] == "" {
		return []string{}
	}
	return split
}

var ErrMissingPermissions = errors.New("missing permissions")

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

	if stat.UserMetadata["Permissions"] == "" {
		return nil, ErrMissingPermissions
	}
	permissions, err := strconv.Atoi(stat.UserMetadata["Permissions"])
	if err != nil {
		return nil, err
	}

	meta := structs.Meta{}

	meta["id"] = stat.UserMetadata["Id"]
	meta["original-name"] = stat.UserMetadata["Original-Name"]
	meta["path"] = stat.Key

	meta["encoding"] = stat.UserMetadata["Encoding"]
	meta["mime-type"] = stat.UserMetadata["Mime-Type"]
	meta["size"] = stat.Size
	meta["checksum"] = stat.ChecksumSHA256
	meta["checksum-type"] = "SHA256"

	meta["owner"] = stat.UserMetadata["Owner"]
	meta["uploaded-by"] = stat.UserMetadata["Uploaded-By"]
	meta["permissions"] = entity.FilePermissions(permissions)

	meta["description"] = stat.UserMetadata["Description"]
	meta["categories"] = splitGroupedField(stat.UserMetadata["Categories"])
	meta["tags"] = splitGroupedField(stat.UserMetadata["Tags"])

	meta["status"] = entity.FileStatus(stat.UserMetadata["Status"])

	meta["uploaded-at"] = stat.UserMetadata["Uploaded-At"]
	meta["updated-at"] = stat.UserMetadata["Updated-At"]
	meta["created-at"] = stat.UserMetadata["Created-At"]
	meta["accessed-at"] = stat.UserMetadata["Accessed-At"]

	structuredMetadata, err := entity.NewFileMetadata(meta)
	if err != nil {
		return nil, err
	}

	return structuredMetadata, nil
}
