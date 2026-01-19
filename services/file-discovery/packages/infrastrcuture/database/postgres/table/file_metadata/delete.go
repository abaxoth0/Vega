package filemetadatatable

import (
	"vega_file_discovery/packages/entity"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

func (_ *Manager) SoftDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.FileMetadata, error) {
	panic("SoftDeleteFileMetadata not implemented")
}

func (_ *Manager) HardDeleteFileMetadata(cmd *cqrs.IdTargetedCommandQuery) (*entity.FileMetadata, error) {
	panic("HardDeleteFileMetadata not implemented")
}
