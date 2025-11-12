package miniocommon

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"vega/packages/domain/entity"
	MinIOConnection "vega/packages/infrastructure/object-storage/MinIO/connection"

	"github.com/abaxoth0/Vega/libs/go/packages/structs"
)

var storage = MinIOConnection.Manager

// MinIO-specific metadata type (UserMetadata)
type RawMetadata map[string]string

var metadataFields = []string{
	"id",
	"original-name",
	"path",
	"encoding",
	"mime-type",
	"checksum",
	"checksum-type",
	"uploaded-by",
	"permissions",
	"description",
	"categories",
	"uploaded-at",
	"updated-at",
	"created-at",
	"accessed-at",
	"status",
}

func toString(v any) string {
	switch src := v.(type) {
	case string:
		return src
	case int64:
		return strconv.FormatInt(src, 10)
	case []string:
		return strings.Join(src, ";")
	case time.Time:
		return src.Format(time.RFC3339)
	case entity.FileStatus:
		return src.String()
	default:
		return ""
	}
}

func ConvertToRawMetadata(meta structs.Meta) RawMetadata {
	r := make(RawMetadata)

	for _, field := range metadataFields {
		r[field] = toString(meta[field])
	}

	return r
}

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
