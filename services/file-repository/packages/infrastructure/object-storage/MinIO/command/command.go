package miniocommand

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

	"github.com/abaxoth0/Vega/libs/go/packages/structs"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
)

var Handler FileApplication.CommandHandler = new(defaultCommandHandler)

type defaultCommandHandler struct {
}

var storage = MinIOConnection.Manager

func (h *defaultCommandHandler) preprocessTargetedCommandQuery(
	commandQuery *application.CommandQuery, path string,
) error {
	if !commandQuery.IsInit() {
		application.InitDefaultCommandQuery(commandQuery)
	}
	if err := FileApplication.ValidatePathFormat(path); err != nil {
		return err
	}
	return nil
}

func (h *defaultCommandHandler) preprocessCommandQuery(
	commandQuery *application.CommandQuery,
) (context.Context, context.CancelFunc) {
	if !commandQuery.IsInit() {
		application.InitDefaultCommandQuery(commandQuery)
	}

	ctx, cancel := context.WithTimeout(commandQuery.Context, commandQuery.ContextTimeout)

	return ctx, cancel
}

func (h *defaultCommandHandler) Mkdir(cmd *FileApplication.MkdirCommand) error {
	if err := h.preprocessTargetedCommandQuery(&cmd.CommandQuery, cmd.Path); err != nil {
		return err
	}
	if ok := FileApplication.IsDirectory(cmd.Path); !ok {
		return FileApplication.ErrFileIsNotDirectory
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	_, err := storage.Client.PutObject(ctx, cmd.Bucket, cmd.Path, nil, 0, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) UploadFile(cmd *FileApplication.UploadFileCommand) error {
	if err := h.preprocessTargetedCommandQuery(&cmd.CommandQuery, cmd.Path); err != nil {
		return err
	}
	if cmd.ContentSize <= 0 {
		return errors.New("Content size cannot be equal or less than 0, but got " + strconv.FormatInt(cmd.ContentSize, 10))
	}
	if FileApplication.IsDirectory(cmd.Path) {
		// TODO (FEAT): Allow upload archives as directories (using flag in cmd?)
		return errors.New("Can't upload file as directory")
	}
	if cmd.Content == nil {
		cmd.Content = bytes.NewReader([]byte{})
	}

	ctx, cancel := context.WithTimeout(cmd.Context, cmd.ContextTimeout)
	defer cancel()

	if err := MinIOCommon.IsBucketExist(ctx, cmd.Bucket); err != nil {
		return err
	}

	hasher := sha256.New()
	teedStream := io.TeeReader(cmd.Content, hasher)
	mimeBuffer := make([]byte, 512)

	n, err := teedStream.Read(mimeBuffer)
	if err != nil && err != io.EOF {
		return err
	}

	mimeType := mimetype.Detect(mimeBuffer[:n])

	multiPartStream := io.MultiReader(bytes.NewReader(mimeBuffer[:n]), teedStream)

	meta := entity.FileMetadata{
		UploadedAt:   time.Now(),
		CreatedAt:    time.Now(), // TODO temp
		Status:       entity.ActiveFileStatus,
		OriginalName: path.Base(cmd.Path),
		Path:         cmd.Path,
		Checksum:     hex.EncodeToString(hasher.Sum(nil)), // // This will be correct after full upload
		ChecksumType: "sha256",
		MIMEType:     mimeType.String(),
	}

	_, err = storage.Client.PutObject(ctx, cmd.Bucket, cmd.Path, multiPartStream, cmd.ContentSize, minio.PutObjectOptions{
		UserMetadata: miniocommon.ConvertToRawMetadata(meta.Pack()),
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) UpdateFileMetadata(cmd *FileApplication.UpdateFileMetadataCommand) error {
	if err := h.preprocessTargetedCommandQuery(&cmd.CommandQuery, cmd.Path); err != nil {
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
		Bucket:          bucket,
		Object:          path,
		UserMetadata:    miniocommon.ConvertToRawMetadata(newMetadata),
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
	if err := h.preprocessTargetedCommandQuery(&cmd.CommandQuery, cmd.Path); err != nil {
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
//
// TODO (FEAT): Implement recursive deletion for directories
func (h *defaultCommandHandler) DeleteFiles(cmd *FileApplication.DeleteFilesCommand) error {

	if !cmd.CommandQuery.IsInit() {
		application.InitDefaultCommandQuery(&cmd.CommandQuery)
	}

	for _, path := range cmd.Paths {
		if err := FileApplication.ValidatePathFormat(path); err != nil {
			return err
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

func (h *defaultCommandHandler) MakeBucket(cmd *FileApplication.MakeBucketCommand) error {
	ctx, cancel := h.preprocessCommandQuery(&cmd.CommandQuery)
	defer cancel()

	if err := storage.Client.MakeBucket(ctx, cmd.Name, minio.MakeBucketOptions{}); err != nil {
		return err
	}

	return nil
}

func (h *defaultCommandHandler) DeleteBucket(cmd *FileApplication.DeleteBucketCommand) error {
	ctx, cancel := h.preprocessCommandQuery(&cmd.CommandQuery)
	defer cancel()

	err := storage.Client.RemoveBucketWithOptions(ctx, cmd.Name, minio.RemoveBucketOptions{
		ForceDelete: cmd.Force,
	})
	if err != nil {
		return err
	}

	return nil
}
