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
	GetFileMetadataByPath(query *GetFileByPathQuery) (*entity.FileMetadata, error)
	SearchFilesByOwner(query *SearchFilesByOwnerQuery) ([]*entity.File, error)
}

type CommandHandler interface {
	Mkdir(cmd *MkdirCommand) error
	UploadFile(cmd *UploadFileCommand) error
	UpdateFileContent(cmd *UpdateFileContentCommand) error
	UpdateFileMetadata(cmd *UpdateFileMetadataCommand) error
	DeleteFiles(cmd *DeleteFilesCommand) error
}
