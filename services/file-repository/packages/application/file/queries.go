package fileapplication

import "vega/packages/application"

type GetFileByPathQuery struct {
	Bucket string
	Path   string

	application.CommandQuery
}

type SearchFilesByOwnerQuery struct {
	Owner string
}
