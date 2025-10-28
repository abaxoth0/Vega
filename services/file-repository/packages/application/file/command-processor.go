package fileapplication

import "vega/packages/domain/entity"

type UseCases interface {
	QueryHandler
	CommandHandler
}

type QueryHandler interface {
	GetFileByID(query GetFileByIdQuery) (*entity.File, error)
}

type CommandHandler interface {
	CreateFile(cmd CreateFileCommand) error
	UpdateFile(cmd UpdateFileCommand) error
	DeleteFile(cmd DeleteFileCommand) error
}

