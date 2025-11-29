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
}

type CommandHandler interface {
	Mkdir(cmd *MkdirCommand) error
	UploadFile(cmd *UploadFileCommand) error
	UpdateFileContent(cmd *UpdateFileContentCommand) error
	DeleteFiles(cmd *DeleteFilesCommand) error
	MakeBucket(cmd *MakeBucketCommand) error
	DeleteBucket(cmd *DeleteBucketCommand) error
}
