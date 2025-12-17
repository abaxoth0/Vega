package fileapplication

import (
	"io"

	"github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

type MkdirCommand struct {
	Bucket string
	Path   string

	cqrs.CommandQuery
}

type UploadFileCommand struct {
	Content     io.Reader
	ContentSize int64
	Path        string
	Bucket      string

	cqrs.CommandQuery
}

type UpdateFileContentCommand struct {
	Path       string
	Bucket     string
	NewContent io.Reader
	Size	   int64

	cqrs.CommandQuery
}

type DeleteFilesCommand struct {
	Paths  []string
	Bucket string

	cqrs.CommandQuery
}

type MakeBucketCommand struct {
	Name string

	cqrs.CommandQuery
}

type DeleteBucketCommand struct {
	Name  string
	Force bool

	cqrs.CommandQuery
}
