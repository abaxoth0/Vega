package filemetadatatable

import (
	_ "embed"
	"vega_file_discovery/packages/entity"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

//go:embed sql/get-file-metadata-by-id.sql
var getFileMetadataByIDSql string;

func (_ *Manager) GetFileMetadataByID(cqrsQuery *cqrs.IdTargetedCommandQuery) (*entity.FileMetadata, error) {
	dblog.Logger.Info("Getting file metadata with id = "+cqrsQuery.ID+"...", nil)

	selectQuery := query.New(getFileMetadataByIDSql, cqrsQuery.ID)

	metadata, err := executor.RowFileMetadata(connection.Primary, selectQuery, "none")
	if err != nil {
		return nil, err
	}

	dblog.Logger.Info("Getting file metadata with id = "+cqrsQuery.ID+": OK", nil)

	return metadata, nil;
}
