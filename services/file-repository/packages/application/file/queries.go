package fileapplication

import "vega/packages/application"

type GetFileByNameQuery struct {
	Bucket 			string	`json:"bucket"`
	Path			string	`json:"path"`

	application.CommandQuery
}

type SearchFilesByOwnerQuery struct {
	Owner	string	`json:"owner"`
}

