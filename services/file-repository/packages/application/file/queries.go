package fileapplication

import "vega/packages/application"

type GetFileByNameQuery struct {
	Bucket 			string
	Path			string

	application.CommandQuery
}

type SearchFilesByOwnerQuery struct {
	Owner	string
}

