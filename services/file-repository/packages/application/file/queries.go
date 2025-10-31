package fileapplication

import (
	"context"
	"time"
)

type GetFileByNameQuery struct {
	FileName		string	`json:"id"`
	Bucket 			string	`json:"bucket"`
	Path			string	`json:"path"`
	Context			context.Context
	ContextTimeout 	time.Duration
}

type SearchFilesByOwnerQuery struct {
	Owner	string	`json:"owner"`
}

