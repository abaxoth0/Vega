package fileapplication

import "errors"

func IsDirectory(path string) bool {
	// If file ends with '/' this file is directory
	return path[len(path)-1] == '/'
}

const (
	maxPathLength int = 1024
	maxPathSegmentLength int = 255
)

func ValidatePathFormat(path string) error {
	if len(path) > maxPathLength {
		return errors.New("max path length exceeded")
	}
	if len(path) == 0 {
		return errors.New("empty path")
	}
	if path[0] != '/' {
		return errors.New("invalid path format: path must begin with \"/\"")
	}
	count := 0
	for char := range path {
		if char == '/' {
			count = 0
			continue
		}
		count++
		if count > maxPathSegmentLength {
			return errors.New("max path's segment length exceeded")
		}
	}
	return nil
}

var ErrInvalidPathFormat = errors.New("Invalid path format")

func FileNameFromPath(path string) (string, error) {
	if err := ValidatePathFormat(path); err != nil {
		return "", err
	}
	if len(path) == 1 {
		return "/", nil
	}
	offset := 0
	if IsDirectory(path) {
		offset++
	}

	fileName := ""
	char := byte(0)

	for i := len(path) - 1 - offset; i >= 0; i--  {
		char = path[i]
		if char == '/' {
			return fileName, nil
		}
		fileName = string(char) + fileName
	}

	return "", ErrInvalidPathFormat
}
