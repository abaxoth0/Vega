package entity

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/abaxoth0/Vega/go-libs/packages/structs"
)

type FileStatus string

const (
	ActiveFileStatus        FileStatus = "active"
	DeletedFileStatus       FileStatus = "deleted"
	ArchivedFileStatus      FileStatus = "archived"
	PendingReviewFilestatus FileStatus = "pending review"
)

var fileStatuses = map[FileStatus]bool{
	ActiveFileStatus:        true,
	DeletedFileStatus:       true,
	ArchivedFileStatus:      true,
	PendingReviewFilestatus: true,
}

func (s FileStatus) Validate() error {
	if _, ok := fileStatuses[s]; !ok {
		return fmt.Errorf("file status \"%s\" doesn't exist", s)
	}
	return nil
}

type FileMetadata struct {
	ID           string
	Originalname string
	StoragePath  string

	Encoding     string
	MIMEType     string
	Size         uint64
	Checksum     string
	ChecksumType string

	Owner       string
	UploadedBy  string
	Permissions string

	Description string
	Categories  []string
	Tags        []string
	ParentDir   string

	UploadedAt time.Time
	UpdatedAt  time.Time
	CreatedAt  time.Time
	AccessedAt time.Time

	Status FileStatus
}

func NewFrom(meta structs.Meta) (*FileMetadata, error) {
	fileMeta := new(FileMetadata)

	fileMeta.ID = meta["id"].(string)
	fileMeta.Originalname = meta["original-name"].(string)
	fileMeta.StoragePath = meta["storage-path"].(string)

	fileMeta.Encoding = meta["encoding"].(string)
	fileMeta.MIMEType = meta["mime-type"].(string)
	fileMeta.Size = meta["size"].(uint64)
	fileMeta.Checksum = meta["checksum"].(string)
	fileMeta.ChecksumType = meta["checksum-type"].(string)

	fileMeta.Owner = meta["owner"].(string)
	fileMeta.UploadedBy = meta["uploaded-by"].(string)
	fileMeta.Permissions = meta["permissions"].(string)

	fileMeta.Description = meta["description"].(string)
	fileMeta.Categories = meta["categories"].([]string)
	fileMeta.Tags = meta["tags"].([]string)
	fileMeta.ParentDir = meta["parent-dir"].(string)

	fileMeta.Status = meta["status"].(FileStatus)
	if err := fileMeta.Status.Validate(); err != nil {
		return nil, err
	}

	rawTimestamps := []any{
		meta["uploaded-at"],
		meta["updated-at"],
		meta["created-at"],
		meta["accessed-at"],
	}

	for _, rawTimestamp := range rawTimestamps {
		switch ts := rawTimestamp.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return nil, errors.New("Invalid timestamp layout: expected RFC3339, but got " + ts)
			}
			fileMeta.UploadedAt = t
		case time.Time:
			fileMeta.UploadedAt = ts
		default:
			return nil, fmt.Errorf("Invalid timestamp type: expected string or time.Time, but got %T", ts)
		}
	}

	return fileMeta, nil
}

func (m *FileMetadata) Pack() structs.Meta {
	meta := make(structs.Meta)

	meta["id"] = m.ID
	meta["original-name"] = m.Originalname
	meta["storage-path"] = m.StoragePath

	meta["encoding"] = m.Encoding
	meta["mime-type"] = m.MIMEType
	meta["size"] = m.Size
	meta["checksum"] = m.Checksum
	meta["checksum-type"] = m.ChecksumType

	meta["owner"] = m.Owner
	meta["uploaded-by"] = m.UploadedBy
	meta["permissions"] = m.Permissions

	meta["description"] = m.Description
	meta["categories"] = m.Categories
	meta["tags"] = m.Tags
	meta["parent-dir"] = m.ParentDir

	meta["uploaded-at"] = m.UploadedAt
	meta["updated-at"] = m.UpdatedAt
	meta["created-at"] = m.CreatedAt
	meta["accessed-at"] = m.AccessedAt

	meta["status"] = m.Status

	return meta
}

func (m *FileMetadata) AddTag(tag string) {
	if slices.Contains(m.Tags, tag) {
		return
	}
	m.Tags = append(m.Tags, tag)
}

func (m *FileMetadata) HasTag(tag string) bool {
	return slices.Contains(m.Tags, tag)
}

func (m *FileMetadata) AddCategory(category string) {
	if slices.Contains(m.Tags, category) {
		return
	}
	m.Categories = append(m.Categories, category)
}

func (m *FileMetadata) HasCategory(category string) bool {
	return slices.Contains(m.Tags, category)
}

type File struct {
	Meta    *FileMetadata
	Content []byte
}

type FileStream struct {
	Reader 	io.Reader
	Context context.Context
	Cancel 	context.CancelFunc
}
