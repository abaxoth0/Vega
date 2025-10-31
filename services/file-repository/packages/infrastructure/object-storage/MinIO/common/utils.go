package miniocommon

func IsDirectory(path string) bool {
	// If file ends with '/' this file is directory
	return path[len(path)-1] == '/'
}

