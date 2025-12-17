package fileapplication

import (
	"errors"

	cqrs "github.com/abaxoth0/Vega/libs/go/packages/CQRS"
)

var (
	ErrFileDoesNotExist   = errors.New("requested file doesn't exist")
	ErrBucketDoesNotExist = errors.New("requested bucket doesn't exist")
)

type GetFileByPathQuery struct {
	Bucket string
	Path   string

	cqrs.CommandQuery
}
