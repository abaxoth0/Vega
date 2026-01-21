package executor

import (
	"database/sql"
	"errors"
	"vega_file_discovery/packages/entity"
	dbcommon "vega_file_discovery/packages/infrastrcuture/database/postgres/common"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
)

var ErrNotSoftDeleted = errors.New("Requested resource exists, but it's not soft deleted")

func rowAnyFileMetadata[T *entity.FileMetadata|*entity.DeletedFileMetadata](
	conType connection.Type, q *query.Query, cacheKey string,
) (T, error){
	scan, err := Row(conType, q)
	if err != nil {
		return nil, err;
	}

	metadata := entity.FileMetadata{}

	var (
		uploadedAt sql.NullTime
		updatedAt  sql.NullTime
		createdAt  sql.NullTime
		accessedAt sql.NullTime
		deletedAt  sql.NullTime
		rawFileStatus string
	)

	dests := []any{
		&metadata.ID,
		&metadata.Path,
		&metadata.Bucket,
		&metadata.Size,
		&metadata.MIMEType,
		&metadata.Permissions,
		&rawFileStatus,
		&metadata.Owner,
		&metadata.OriginalName,
		&metadata.Encoding,
		&metadata.Categories,
		&metadata.Tags,
		&metadata.Checksum,
		&metadata.UploadedBy,
		&uploadedAt,
		&accessedAt,
		&updatedAt,
		&createdAt,
		&metadata.Description,
	}

	var result T

	switch any(result).(type) {
	case *entity.DeletedFileMetadata:
		dests = append(dests, &deletedAt)
	}

	if err := scan(dests...); err != nil {
		return nil, err;
	}
	if uploadedAt.Valid {
		metadata.UploadedAt = uploadedAt.Time
	}
	if accessedAt.Valid {
		metadata.AccessedAt = accessedAt.Time
	}
	if updatedAt.Valid {
		metadata.UpdatedAt = updatedAt.Time
	}
	if createdAt.Valid {
		metadata.CreatedAt = createdAt.Time
	}

	switch m := any(result).(type) {
	case *entity.DeletedFileMetadata:
		if m == nil {
			m = new(entity.DeletedFileMetadata)
		}
		if deletedAt.Valid {
			m.DeletedAt = deletedAt.Time
		} else {
			return nil, ErrNotSoftDeleted
		}
		m.FileMetadata = metadata
	}

	var e error
	metadata.Status, e = dbcommon.ParseFileStatus(rawFileStatus)
	if e != nil {
		dblog.Logger.Error(
			"Failed to get file metadata " + metadata.ID,
			e.Error(),
			nil,
		)
		// The error may occur only if raw file status in DB is invalid
		// TODO: try to resolve this problem (apply previous status or set status to pending)
		return nil, errs.StatusInternalServerError
	}

	return result, nil
}

func RowFileMetadata(conType connection.Type, q *query.Query, cacheKey string) (*entity.FileMetadata, error){
	return rowAnyFileMetadata[*entity.FileMetadata](conType, q, cacheKey)
}

func RowSoftDeletedFileMetadata(
	conType connection.Type, q *query.Query, cacheKey string,
) (*entity.DeletedFileMetadata, error){
	return rowAnyFileMetadata[*entity.DeletedFileMetadata](conType, q, cacheKey)
}
