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

type UpdateFileCommand struct {

}

type DeleteFilesCommand struct {
	Paths	[]string
	Bucket 	string

	application.CommandQuery
}

