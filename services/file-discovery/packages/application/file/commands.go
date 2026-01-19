package fileapplication

import (
	"vega_file_discovery/packages/entity"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

type CreateFileMetadataCmd struct {
	Metadata *entity.FileMetadata

	cqrs.CommandQuery
}

type UpdateFileMetadataCmd struct {
	Upd  *entity.UpdatableFileMetadata

	cqrs.IdTargetedCommandQuery
}
