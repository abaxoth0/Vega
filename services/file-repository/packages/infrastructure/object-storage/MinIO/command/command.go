package miniocommand

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
	"vega/packages/application"
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
	MinIOCommon "vega/packages/infrastructure/object-storage/MinIO/common"
	miniocommon "vega/packages/infrastructure/object-storage/MinIO/common"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/abaxoth0/Vega/go-libs/packages/structs"
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
	if err := FileApplication.ValidatePathFormat(cmd.Path); err != nil {
		return err
	}
	if FileApplication.IsDirectory(cmd.Path) {
		// TODO (FEAT): Allow upload archives as directories (using flag in cmd?)
		return errors.New("Can't upload file as directory")
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	if err := MinIOCommon.IsBucketExist(ctx, cmd.Bucket); err != nil {
		return err
	}

	meta := entity.FileMetadata{
		UploadedAt: time.Now(),
		CreatedAt: time.Now(), // TODO temp
		Status: entity.ActiveFileStatus,
		Size: cmd.ContentSize,
		OriginalName: path.Base(cmd.Path),
		Path: cmd.Path,
	}

	_, err := storage.Client.PutObject(ctx, cmd.Bucket, cmd.Path, cmd.Content, cmd.ContentSize, minio.PutObjectOptions{
		UserMetadata: miniocommon.ConvertToRawMetadata(meta.Pack()),
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) UpdateFileMetadata(cmd *FileApplication.UpdateFileMetadataCommand) error {
	if !cmd.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&cmd.CommandQuery)
	}
	if err := FileApplication.ValidatePathFormat(cmd.Path); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()


	info, err := storage.Client.StatObject(ctx, cmd.Bucket, cmd.Path, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	if info.Size < entity.SmallFileSizeThreshold {
		return h.copyWithNewMetadata(ctx, cmd.Bucket, cmd.Path, cmd.NewMetadata)
	}

	return errors.New("bigus")
}

func (h *defaultCommandHandler) copyWithNewMetadata(
	ctx context.Context,
	bucket string,
	path string,
	newMetadata structs.Meta,
) error {
	src := minio.CopySrcOptions{
		Bucket: bucket,
		Object: path,
	}
	dest := minio.CopyDestOptions{
		Bucket: bucket,
		Object: path,
		UserMetadata: miniocommon.ConvertToRawMetadata(newMetadata),
		ReplaceMetadata: true,
	}

	_, err := storage.Client.CopyObject(ctx, dest, src)
	if err != nil {
		return err
	}
	return nil
}

func (h *defaultCommandHandler) fullReplace(ctx context.Context, bucket string, path string, content []byte) error {
	size := int64(len(content))

	_, err := storage.Client.PutObject(ctx, bucket, path, bytes.NewReader(content), size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) UpdateFileContent(cmd *FileApplication.UpdateFileContentCommand) error {
	if !cmd.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&cmd.CommandQuery)
	}
	if err := FileApplication.ValidatePathFormat(cmd.Path); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	if err := h.fullReplace(ctx, cmd.Bucket, cmd.Path, cmd.NewContent); err != nil {
		return err
	}

	return nil
}

// TODO (FEAT?): By default MinIO doesn't consider situation when you trying to delete non-existing file as error.
// Which is reasonable decision, since this file doesn't exist at the end - operation can be considered successful.
// But in some cases it may be important for end user to know, does this file even existed or not? So maybe add
// possibility for users to decide - should file existance be checked before deletion or not? But this can be used
// for possible attacks, like DoS... so - does it even worth this? I don't know, i can't really imagine situations
// when this functional will be really needed... So maybe leave it as it is works now? Again - i don't know...
func (h *defaultCommandHandler) DeleteFiles(cmd *FileApplication.DeleteFilesCommand) error {
	if !cmd.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&cmd.CommandQuery)
	}

	for _, path := range cmd.Paths {
		if err := FileApplication.ValidatePathFormat(path); err != nil {
			return err
		}
		// TODO (FEAT): Implement recursive deletion for directories
		if FileApplication.IsDirectory(path) {
			return errors.New("can't delete directory: " + path)
		}
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	if err := MinIOCommon.IsBucketExist(ctx, cmd.Bucket); err != nil {
		return err
	}

	if len(cmd.Paths) == 1 {
		err := storage.Client.RemoveObject(ctx, cmd.Bucket, cmd.Paths[0], minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	objectsCh := make(chan minio.ObjectInfo, len(cmd.Paths))

	for _, path := range cmd.Paths {
		objectsCh <- minio.ObjectInfo{Key: path}
	}
	close(objectsCh)

	errorCh := storage.Client.RemoveObjects(ctx, cmd.Bucket, objectsCh, minio.RemoveObjectsOptions{})

	var errors []string
	for err := range errorCh {
		if err.Err != nil {
			errors = append(errors, fmt.Sprintf("Failed to delete %s: %v", err.ObjectName, err.Err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("deletion errors: %s", strings.Join(errors, ";"))
	}

	return nil
}
