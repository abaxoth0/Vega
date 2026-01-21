package filemetadatatable

import (
	_ "embed"
	"time"
	"vega_file_discovery/packages/entity"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/executor"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

//go:embed sql/soft-delete-file-metadata.sql
var softDeleteFileMetadataSql string;

//go:embed sql/hard-delete-file-metadata.sql
var hardDeleteFileMetadataSql string;

func (m *Manager) SoftDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.DeletedFileMetadata, error) {
	dblog.Logger.Info("Soft deleting file metadata "+cmd.ID+"...", nil)

	metadata, err := m.GetFileMetadataByID(&cqrs.IdTargetedCommandQuery{
		ID: cmd.ID,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if err := executor.Exec(
		connection.Primary,
		query.New(softDeleteFileMetadataSql, now, cmd.ID),
	); err != nil {
		return nil, err
	}

	dblog.Logger.Info("Soft deleting file metadata "+cmd.ID+":OK", nil)

	return &entity.DeletedFileMetadata{
		FileMetadata: *metadata,
		DeletedAt: now,
	}, nil
}

func (m *Manager) HardDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.DeletedFileMetadata, error) {
	dblog.Logger.Info("Hard deleting file metadata "+cmd.ID+"...", nil)

	deletedMetadata, err := m.getSoftDeletedFileMetadataByID(&cqrs.IdTargetedCommandQuery{
		ID: cmd.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := executor.Exec(
		connection.Primary,
		query.New(hardDeleteFileMetadataSql, cmd.ID),
	); err != nil {
		return nil, err
	}

	dblog.Logger.Info("Hard deleting file metadata "+cmd.ID+":OK", nil)

	return deletedMetadata, nil
}
