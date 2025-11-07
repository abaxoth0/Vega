package fileapplication

import (
	"io"
	"vega/packages/application"
	"vega/packages/domain/entity"
)

type UploadFileCommand struct {
	FileMeta 	*entity.FileMetadata
	Content		io.Reader
	ContentSize	int64
	Path		string
	Bucket 		string

	application.CommandQuery
}

type UpdateFileMetadataCommand struct {
	Path 		string
	Bucket		string
	NewMetadata	*entity.FileMetadata

	application.CommandQuery
}

type UpdateFileContentCommand struct {
	Path 		string
	Bucket 		string
	NewContent 	[]byte

	application.CommandQuery
}

type DeleteFilesCommand struct {
	Paths	[]string
	Bucket 	string

	application.CommandQuery
}

