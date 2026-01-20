package executor

import (
	"database/sql"
	"vega_file_discovery/packages/entity"
	dbcommon "vega_file_discovery/packages/infrastrcuture/database/postgres/common"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
)

func RowFileMetadata(conType connection.Type, q *query.Query, cacheKey string) (*entity.FileMetadata, error){
	scan, err := Row(conType, q)
	if err != nil {
		return nil, err;
	}

	metadata := new(entity.FileMetadata)

	var (
		uploadedAt sql.NullTime
		updatedAt  sql.NullTime
		createdAt  sql.NullTime
		accessedAt sql.NullTime
		rawFileStatus string
	)

	if err := scan(
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
	); err != nil {
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

	return metadata, nil
}
