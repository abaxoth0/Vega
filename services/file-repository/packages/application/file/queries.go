package fileapplication

import "vega/packages/application"

type GetFileByPathQuery struct {
	Bucket string
	Path   string

	application.CommandQuery
}
