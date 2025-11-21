package fileapplication

import (
	"io"
	"vega/packages/application"
	"vega/packages/domain/entity"

	"github.com/abaxoth0/Vega/libs/go/packages/structs"
)

type MkdirCommand struct {
	Bucket 	string
	Path 	string

	application.CommandQuery
}

type UploadFileCommand struct {
	FileMeta    *entity.FileMetadata
	Content     io.Reader
	ContentSize int64
	Path        string
	Bucket      string

	application.CommandQuery
}

type UpdateFileMetadataCommand struct {
	Path        string
	Bucket      string
	NewMetadata structs.Meta

	application.CommandQuery
}

type UpdateFileContentCommand struct {
	Path       string
	Bucket     string
	NewContent []byte

	application.CommandQuery
}

type DeleteFilesCommand struct {
	Paths  []string
	Bucket string

	application.CommandQuery
}

type BucketCommand struct {
	Name string

	application.CommandQuery
}
