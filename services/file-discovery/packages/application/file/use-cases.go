package fileapplication

import (
	"vega_file_discovery/packages/entity"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

type UseCases interface {
	QueryHandler
	CommandHandler
}

type QueryHandler interface {
	GetFileMetadataByID(query *cqrs.IdTargetedCommandQuery) (*entity.FileMetadata, error)
}

type CommandHandler interface {
	CreateFileMetadata(cmd *CreateFileMetadataCmd) (id string, err error)
	UpdateFileMetadata(cmd *UpdateFileMetadataCmd) error
	SoftDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.DeletedFileMetadata, error)
	HardDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.DeletedFileMetadata, error)
}
