package entity

import (
	"context"
	"io"
)

const (
	SmallFileSizeThreshold = 10 * 1024 * 1024 // 10 MB
)

type File struct {
	Content []byte
}

type FileStream struct {
	Content io.Reader
	Size    int64
	Context context.Context
	Cancel  context.CancelFunc
}
