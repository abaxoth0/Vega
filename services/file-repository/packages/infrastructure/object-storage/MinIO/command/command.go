package miniocommand

import (
	"context"
	"errors"
	"strconv"
	"vega/packages/application"
	FileApplication "vega/packages/application/file"
	MinIOCommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/minio/minio-go/v7"
)

var Handler FileApplication.CommandHandler = new(defaultCommandHandler)

type defaultCommandHandler struct {

}

var storage = MinIOConnection.Manager

func (h *defaultCommandHandler) UploadFile(cmd *FileApplication.UploadFileCommand) error {
	if !cmd.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&cmd.CommandQuery)
	}

	if cmd.ContentSize <= 0 {
		return errors.New("Content size cannot be equal or less than 0, but got " + strconv.FormatInt(cmd.ContentSize, 10))
	}
	if FileApplication.IsDirectory(cmd.Path) {
		// TODO Allow upload archives as directories (using flag in cmd?)
		return errors.New("Can't upload file as directory")
	}
	if err := FileApplication.ValidatePathFormat(cmd.Path); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	if err := MinIOCommon.IsBucketExist(ctx, cmd.Bucket); err != nil {
		return err
	}
	_, err := storage.Client.PutObject(ctx, cmd.Bucket, cmd.Path, cmd.Content, cmd.ContentSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) UpdateFile(cmd *FileApplication.UpdateFileCommand) error {
	return nil
}

func (h *defaultCommandHandler) DeleteFiles(cmd *FileApplication.DeleteFilesCommand) error {
	return nil
}

