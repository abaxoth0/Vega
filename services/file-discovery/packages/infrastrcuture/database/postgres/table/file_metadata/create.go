package filemetadatatable

import (
	_ "embed"
	"time"
	fileapplication "vega_file_discovery/packages/application/file"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	"github.com/google/uuid"
)

//go:embed sql/insert-file-metadata.sql
var insertFileMetadataSql string;

func (_ *Manager) CreateFileMetadata(cmd *fileapplication.CreateFileMetadataCmd) (string, error) {
	dblog.Logger.Info("Creating new file metadata", nil)

	var (
		uploadedAt time.Time
		createdAt  time.Time
	)
	if cmd.Metadata.UploadedAt.IsZero() {
		uploadedAt = time.Now()
	}
	if cmd.Metadata.CreatedAt.IsZero() {
		createdAt = time.Now()
	}

	id := uuid.New()
	insertQuery := query.New(
		insertFileMetadataSql,
		id,
		cmd.Metadata.Path,
		cmd.Metadata.Bucket,
		cmd.Metadata.Size,
		cmd.Metadata.MIMEType,
		cmd.Metadata.Permissions,
		convertFileStatus(cmd.Metadata.Status),
		cmd.Metadata.Owner,
		cmd.Metadata.OriginalName,
		cmd.Metadata.Encoding,
		query.Nullif(cmd.Metadata.Categories),
		query.Nullif(cmd.Metadata.Tags),
		cmd.Metadata.Checksum,
		cmd.Metadata.UploadedBy,
		uploadedAt,
		createdAt,
		query.Nullif(cmd.Metadata.Description),
	)

	if err := executor.Exec(connection.Primary, insertQuery); err != nil {
		return "", err
	}

	dblog.Logger.Info("Creating new file metadata: OK", nil)

	return id.String(), nil
}
