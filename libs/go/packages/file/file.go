package file

import "errors"

var ErrFileIsNotDirectory = errors.New("file is not directory")

func IsDirectory(path string) bool {
	// If file ends with '/' this file is directory
	return path[len(path)-1] == '/'
}

const (
	maxPathLength        int = 1024
	maxPathSegmentLength int = 255
)

var ErrMaxPathLengthExceeded = errors.New("max path length exceeded")
var ErrEmptyPath = errors.New("empty path")
var ErrInvalidPathFormat = errors.New("invalid path format: path must begin with \"/\"")
var ErrMaxPathSegmentLengthExceeded = errors.New("max path's segment length exceeded")

func ValidatePathFormat(path string) error {
	if len(path) > maxPathLength {
		return ErrMaxPathLengthExceeded
	}
	if len(path) == 0 {
		return ErrEmptyPath
	}
	if path[0] != '/' {
		return ErrInvalidPathFormat
	}
	count := 0
	for char := range path {
		if char == '/' {
			count = 0
			continue
		}
		count++
		if count > maxPathSegmentLength {
			return ErrMaxPathSegmentLengthExceeded
		}
	}
	return nil
}
