package fileapplication

import (
	"io"
	"vega/packages/application"
)

type MkdirCommand struct {
	Bucket string
	Path   string

	application.CommandQuery
}

type UploadFileCommand struct {
	Content     io.Reader
	ContentSize int64
	Path        string
	Bucket      string

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

type MakeBucketCommand struct {
	Name string

	application.CommandQuery
}

type DeleteBucketCommand struct {
	Name  string
	Force bool

	application.CommandQuery
}
