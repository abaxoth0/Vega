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
	Bucket 		*string
	Path   		*string
	Encoding 	*string
	Owner       *string
	Permissions *entity.FilePermissions
	Description *string
	Categories  []string
	Tags        []string
	Status 		*entity.FileStatus

	cqrs.IdTargetedCommandQuery
}
