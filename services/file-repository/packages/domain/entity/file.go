package entity

import (
	"context"
	"io"
)

type FileStream struct {
	Content io.Reader
	Size    int64
	Context context.Context
	Cancel  context.CancelFunc
}
