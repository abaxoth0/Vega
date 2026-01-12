package filemetatable

import (
	"vega_file_discovery/packages/entity"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

func (_ *Manager) GetFileMetadataByID(query *cqrs.IdTargetedCommandQuery) (*entity.FileMetadata, error) {
	panic("GetFileMetadataByID not implemented")
}
