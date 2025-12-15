package fileapplication

import (
	"errors"
	"vega_file_repository/packages/application"
)

var (
	ErrFileDoesNotExist   = errors.New("requested file doesn't exist")
	ErrBucketDoesNotExist = errors.New("requested bucket doesn't exist")
)

type GetFileByPathQuery struct {
	Bucket string
	Path   string

	application.CommandQuery
}
