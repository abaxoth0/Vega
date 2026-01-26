package filemetadatatable

import (
	"slices"
	"strconv"
	"time"
	fileapplication "vega_file_discovery/packages/application/file"
	"vega_file_discovery/packages/entity"
	dbcommon "vega_file_discovery/packages/infrastrcuture/database/postgres/common"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	common "github.com/abaxoth0/Vega/libs/go/packages"
	"github.com/abaxoth0/Vega/libs/go/packages/file"
	"github.com/google/uuid"
)

func posParam(field string, query *query.Query) string {
	return field+" = $"+strconv.Itoa(len(query.Args))
}

func addParam(query *query.Query, field string, value any) {
	query.Args = append(query.Args, value)
	if len(query.Args) > 1 {
		query.SQL += ","
	}
	query.SQL += posParam(field, query)
}

func buildUpdateFileMetadataQuery(
	src *entity.UpdatableFileMetadata, upd *fileapplication.UpdateFileMetadataCmd,
) (*query.Query, error) {
	// TODO Need to update file MIME type, size and checksum if it's content changed
	query := query.New("")

	if upd.Bucket != nil && src.Bucket != *upd.Bucket {
		if err := uuid.Validate(*upd.Bucket); err != nil {
			return nil, err
		}
		addParam(query, "bucket", *upd.Bucket)
	}

	if upd.Path != nil && src.Path != *upd.Path {
		if err := file.ValidatePathFormat(*upd.Path); err != nil {
			return nil, err
		}
		addParam(query, "path", *upd.Path)
	}

	if upd.Encoding != nil && src.Encoding != *upd.Encoding {
		addParam(query, "encoding", common.Ternary(*upd.Encoding == "", "NULL", *upd.Encoding))
	}

	if upd.Owner != nil && src.Owner != *upd.Owner {
		if err := uuid.Validate(*upd.Owner); err != nil {
			return nil, err
		}
		addParam(query, "owner", *upd.Owner)
	}

	if upd.Permissions != nil && src.Permissions != *upd.Permissions {
		if *upd.Permissions == 0 {
			return nil, entity.ErrInvalidFilePermissions
		}
		if upd.Permissions.IsEmpty() {
			return nil, entity.ErrEmptyFilePermissions
		}
		addParam(query, "permissions", upd.Permissions.RawString())
	}

	if upd.Description != nil && src.Description != *upd.Description {
		addParam(query, "description", common.Ternary(*upd.Description == "", "NULL", *upd.Description))
	}

	if slices.Equal(src.Categories, upd.Categories) {
		addParam(query, "categories", dbcommon.SqlArrayFromSlice(upd.Categories))
	}

	if slices.Equal(src.Tags, upd.Tags) {
		addParam(query, "tags", dbcommon.SqlArrayFromSlice(upd.Tags))
	}

	if upd.Status != nil && src.Status != *upd.Status {
		if err := upd.Status.Validate(); err != nil {
			return nil, err
		}
		addParam(query, "status", dbcommon.ConvertFileStatus(*upd.Status))
	}

	updatedAt := time.Now()
	addParam(query, "updated_at", updatedAt)

	query.Args = append(query.Args, upd.ID)
	query.SQL = "UPDATE file_metadata SET "+query.SQL+" WHERE "+posParam("id", query)+";"

	return query, nil
}

func (m *Manager) UpdateFileMetadata(cmd *fileapplication.UpdateFileMetadataCmd) error {
	dblog.Logger.Info("Updating file metadata with id = "+cmd.ID+"...", nil)

	metadata, err := m.GetFileMetadataByID(&cmd.IdTargetedCommandQuery)
	if err != nil {
		return err
	}

	updQuery, err := buildUpdateFileMetadataQuery(&metadata.UpdatableFileMetadata, cmd)
	if err != nil {
		return err
	}

	if err := executor.Exec(connection.Primary, updQuery); err != nil {
		return err
	}

	dblog.Logger.Info("Updating file metadata with id = "+cmd.ID+": OK", nil)

	return nil
}
