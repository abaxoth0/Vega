package fileapplication

import (
	"vega/packages/domain/entity"
)

type UseCases interface {
	QueryHandler
	CommandHandler
}

type QueryHandler interface {
	GetFileByPath(query *GetFileByPathQuery) (*entity.FileStream, error)
	SearchFilesByOwner(query *SearchFilesByOwnerQuery) ([]*entity.File, error)
}

type CommandHandler interface {
	UploadFile(cmd *UploadFileCommand) error
	UpdateFile(cmd *UpdateFileCommand) error
	DeleteFiles(cmd *DeleteFilesCommand) error
}

