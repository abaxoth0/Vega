package miniocommon

import (
	"context"
	"errors"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"
)

var storage = MinIOConnection.Manager

var ErrBucketDoesntExist = errors.New("Bucket doesn't exist")

func IsBucketExist(ctx context.Context, bucket string) error {
	ok, err := storage.Client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !ok {
		return ErrBucketDoesntExist
	}
	return nil
}

