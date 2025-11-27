package entity

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
	"github.com/abaxoth0/Vega/libs/go/packages/structs"
)

const (
	SmallFileSizeThreshold = 10 * 1024 * 1024 // 10 MB
)

type FileStatus string

const (
	ActiveFileStatus        FileStatus = "active"
	DeletedFileStatus       FileStatus = "deleted"
	ArchivedFileStatus      FileStatus = "archived"
	PendingReviewFilestatus FileStatus = "pending review"
)

const (
	DefaultFileStatus FileStatus = ActiveFileStatus
	DefaultChecksumType string = "SHA256"
)

var DefaultChecksumHasher = sha256.New()

var fileStatuses = map[FileStatus]bool{
	ActiveFileStatus:        true,
	DeletedFileStatus:       true,
	ArchivedFileStatus:      true,
	PendingReviewFilestatus: true,
}

func (s FileStatus) Validate() error {
	if s == "" {
		return errors.New("file status is empty")
	}
	if _, ok := fileStatuses[s]; !ok {
		return fmt.Errorf("file status \"%s\" doesn't exist", s)
	}
	return nil
}

func (s FileStatus) String() string {
	return string(s)
}

type FilePermissionGroup uint16

const (
	FilePermissionGroupSize int = 3
	// Amount of FilePermissionGroup in FilePermissions
	FilePermissionGroupsAmount int = 3
)

const (
	DeleteFilePermission FilePermissionGroup = 1 << iota
	UpdateFilePermission
	ReadFilePermission
)

const (
	ReadFilePermissionChar 	 rune = 'r'
	UpdateFilePermissionChar rune = 'u'
	DeleteFilePermissionChar rune = 'd'
	NoFilePermissionChar 	 rune = '-'
)

var permissionMap = map[FilePermissionGroup]rune{
	ReadFilePermission:   ReadFilePermissionChar,
	UpdateFilePermission: UpdateFilePermissionChar,
	DeleteFilePermission: DeleteFilePermissionChar,
}

func (p FilePermissionGroup) String() string {
    str := make([]rune, FilePermissionGroupSize)

    // Process from most significant to least significant bit
	for i := range FilePermissionGroupSize {
        permission := FilePermissionGroup(1 << (FilePermissionGroupSize - 1 - i))

        if p&permission != 0 {
            char, ok := permissionMap[permission]
            if !ok {
                panic(fmt.Sprintf(
                    "invalid permission char \"%s\" for 0b%s",
                    string(char), strconv.FormatInt(int64(permission), 2),
                ))
            }
            str[i] = char
        } else {
            str[i] = NoFilePermissionChar
        }
    }

    return string(str)
}

// File permissions works similar to Linux file permissions.
// And mechanism of how it works can be extended indefinitely, proof:
//
// Consider a set of distinct powers of two:
//
//  { 2^0, 2^1, 2^2, …, 2^(N−1) }
//
// In binary, each number has a single 1 in a unique position. For example:
//
//  1 = 2^0 = 001
//  2 = 2^1 = 010
//  4 = 2^2 = 100
//  ...and so on
//
// Any subset of these numbers corresponds to a binary number where the 1-bits
// indicate which elements are included. For instance:
//
//  { 1, 4 } -> 1 + 4 = 5 (0101)
//
// Because each power of two occupies a unique bit position, every subset produces a unique binary number.
// Hence, all subset sums are distinct.
type FilePermissions uint32

var alignmentBit = func() FilePermissionGroup {
	return FilePermissionGroup(1 << (FilePermissionGroupSize*(FilePermissionGroupSize-1)))
}()

func NewFilePermissions(owner, shared, other FilePermissionGroup) FilePermissions {
	offset := FilePermissionGroup(FilePermissionGroupSize)
	// All bits must be aligned in groups by 3, so to keep this alignment need to add special bit,
	// which will come right after the most left block, by doing so, number will have this format:
	//
	// 0b1<owner-group><shared-group><other-group>
	//   ↑
	//   alignment bit
	return FilePermissions(((alignmentBit | (owner << offset) | shared) << offset) | other)
}

var emptyFilePermissions = NewFilePermissions(0,0,0)

func (p FilePermissions) IsEmpty() bool {
	return p == emptyFilePermissions
}

// Binary number wich is equal to N ones that comes sequentially, N is FilePermissionGroupSize.
// e.g. if FilePermissionGroupSize = 3, then this number will be equal to 0b111.
var filePermissionsGroupBits FilePermissions = func() FilePermissions {
	var r FilePermissions = 1
	for i := 1; i < FilePermissionGroupSize; i++ {
		r |= (1 << i)
	}
	return r
}()

const (
	otherGroupRightOffset int = FilePermissionGroupSize * iota
	sharedGroupRightOffset
	ownerGroupRightOffset
)

func (p FilePermissions) getGroupWithRightOffset(offset int) FilePermissionGroup {
	return FilePermissionGroup(p >> offset & filePermissionsGroupBits)
}

func (p FilePermissions) GetOwnerPermissions() FilePermissionGroup {
	return p.getGroupWithRightOffset(ownerGroupRightOffset)
}

func (p FilePermissions) GetSharedPermissions() FilePermissionGroup {
	return p.getGroupWithRightOffset(sharedGroupRightOffset)
}

func (p FilePermissions) GetOtherPermissions() FilePermissionGroup {
	return p.getGroupWithRightOffset(otherGroupRightOffset)
}

func (p FilePermissions) String() string {
	str := make([]string, FilePermissionGroupsAmount)

	for i := range FilePermissionGroupsAmount {
		str[FilePermissionGroupsAmount-i-1] = p.getGroupWithRightOffset(FilePermissionGroupSize*i).String()
	}

	return strings.Join(str, "")
}

var ErrInvalidFilePermissionGroupLength = errors.New("invalid file permission group length")
var ErrInvalidFilePermissionGroupFormat = errors.New("invalid file permission group format")

func ParseFilePermissionGroup(g string) (FilePermissionGroup, error) {
	if len(g) != FilePermissionGroupSize {
		return 0, ErrInvalidFilePermissionGroupLength
	}

	group := []rune(g)

	fpGroup := FilePermissionGroup(0)

	// TODO probably this should be optimised

	read   := rune(group[0])
	update := rune(group[1])
	del    := rune(group[2])

	if read != ReadFilePermissionChar && read != NoFilePermissionChar {
		return 0, ErrInvalidFilePermissionGroupFormat
	}
	if update != UpdateFilePermissionChar && update != NoFilePermissionChar {
		return 0, ErrInvalidFilePermissionGroupFormat
	}
	if del != DeleteFilePermissionChar && del != NoFilePermissionChar {
		return 0, ErrInvalidFilePermissionGroupFormat
	}

	if read == ReadFilePermissionChar {
		fpGroup |= ReadFilePermission
	}
	if update == UpdateFilePermissionChar {
		fpGroup |= UpdateFilePermission
	}
	if del == DeleteFilePermissionChar {
		fpGroup |= DeleteFilePermission
	}

	return fpGroup, nil
}

var ErrInvalidFilePermissionsLength = errors.New("invalid file permissions length")

func ParseFilePermissions(permissions string) (FilePermissions, error) {
	if len(permissions) != FilePermissionGroupsAmount*FilePermissionGroupSize {
		return 0, ErrInvalidFilePermissionsLength
	}

	rawGroups := make([]string, FilePermissionGroupsAmount)
	rawGroupIdx := 0
	offset := 0

	for i := range permissions {
		if (i+1)%FilePermissionGroupSize == 0 {
			rawGroups[rawGroupIdx] = permissions[offset:FilePermissionGroupSize*(rawGroupIdx+1)]
			offset += FilePermissionGroupSize
			rawGroupIdx++
		}
	}

	if rawGroupIdx != FilePermissionGroupsAmount {
		log.Printf(
			"[ ERROR ] Mistmatch of file permissions groups amount (%d) and amount of parsed groups (%d)\n",
			FilePermissionGroupsAmount, rawGroupIdx+1,
		)
		return 0, errs.StatusInternalServerError
	}

	groups := []FilePermissionGroup{}

	for _, rawGroup := range rawGroups {
		group, err := ParseFilePermissionGroup(rawGroup)
		if err != nil {
			return 0, err
		}
		groups = append(groups, group)
	}

	return NewFilePermissions(groups[0], groups[1], groups[2]), nil
}

type FileMetadata struct {
	ID           string
	OriginalName string
	Path         string

	Encoding     string
	MIMEType     string
	Size         int64
	Checksum     string
	ChecksumType string

	Owner       string
	UploadedBy  string
	Permissions FilePermissions

	Description string
	Categories  []string
	Tags        []string

	UploadedAt time.Time
	UpdatedAt  time.Time
	CreatedAt  time.Time
	AccessedAt time.Time

	Status FileStatus
}

func NewFileMetadata(meta structs.Meta) (*FileMetadata, error) {
	fileMeta := new(FileMetadata)

	fileMeta.ID = meta["id"].(string)
	fileMeta.OriginalName = meta["original-name"].(string)
	fileMeta.Path, _ = meta["storage-path"].(string)

	fileMeta.Encoding = meta["encoding"].(string)
	fileMeta.MIMEType = meta["mime-type"].(string)
	fileMeta.Size = meta["size"].(int64)
	fileMeta.Checksum = meta["checksum"].(string)
	fileMeta.ChecksumType = meta["checksum-type"].(string)

	fileMeta.Owner = meta["owner"].(string)
	fileMeta.UploadedBy = meta["uploaded-by"].(string)
	fileMeta.Permissions = meta["permissions"].(FilePermissions)

	fileMeta.Description = meta["description"].(string)
	fileMeta.Categories = meta["categories"].([]string)
	fileMeta.Tags = meta["tags"].([]string) // TODO: refactor this from []string to map[string]string?

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
	meta["original-name"] = m.OriginalName
	meta["path"] = m.Path

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
	if slices.Contains(m.Categories, category) {
		return
	}
	m.Categories = append(m.Categories, category)
}

func (m *FileMetadata) HasCategory(category string) bool {
	return slices.Contains(m.Categories, category)
}

type File struct {
	Meta    *FileMetadata
	Content []byte
}

type FileStream struct {
	Content io.Reader
	Size    int64
	Context context.Context
	Cancel  context.CancelFunc
}
