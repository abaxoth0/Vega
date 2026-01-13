package entity

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
)

type FileStatus int

const (
	ActiveFileStatus FileStatus = 1 + iota
	ArchivedFileStatus
	PendingFilestatus
)

const DefaultFileStatus = ActiveFileStatus

const (
	activeFileStatusString   string = "active"
	archivedFileStatusString string = "archived"
	pendingFilestatusString  string = "pending"
)

var fileStatusToStrMap = map[FileStatus]string{
	ActiveFileStatus:   activeFileStatusString,
	ArchivedFileStatus: archivedFileStatusString,
	PendingFilestatus:  pendingFilestatusString,
}

var strToFileStatusMap = map[string]FileStatus{
	activeFileStatusString:   ActiveFileStatus,
	archivedFileStatusString: ArchivedFileStatus,
	pendingFilestatusString:  PendingFilestatus,
}

var ErrInvalidFileStatusString = errors.New("invalid string representation of file status")

func ParseFileStatus(status string) (FileStatus, error) {
	r, ok := strToFileStatusMap[status]
	if !ok {
		return 0, ErrInvalidFileStatusString
	}
	return r, nil
}

func (s FileStatus) String() string {
	str := fileStatusToStrMap[s]
	if str == "" {
		return ""
	}
	return str
}

var ErrFileStatuZero = errors.New("file status is zero")

func (s FileStatus) Validate() error {
	if s == 0 {
		return ErrFileStatuZero
	}
	if _, ok := fileStatusToStrMap[s]; !ok {
		return fmt.Errorf("file status \"%s\" doesn't exist", s)
	}
	return nil
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
	ReadFilePermissionChar   rune = 'r'
	UpdateFilePermissionChar rune = 'u'
	DeleteFilePermissionChar rune = 'd'
	NoFilePermissionChar     rune = '-'
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
//	{ 2^0, 2^1, 2^2, …, 2^(N−1) }
//
// In binary, each number has a single 1 in a unique position. For example:
//
//	1 = 2^0 = 001
//	2 = 2^1 = 010
//	4 = 2^2 = 100
//	...and so on
//
// Any subset of these numbers corresponds to a binary number where the 1-bits
// indicate which elements are included. For instance:
//
//	{ 1, 4 } -> 1 + 4 = 5 (0101)
//
// Because each power of two occupies a unique bit position, every subset produces a unique binary number.
// Hence, all subset sums are distinct.
type FilePermissions uint16

var alignmentBit = func() FilePermissionGroup {
	return FilePermissionGroup(1 << (FilePermissionGroupSize * (FilePermissionGroupSize - 1)))
}()

func NewFilePermissions(owner, shared, other FilePermissionGroup) FilePermissions {
	offset := FilePermissionGroup(FilePermissionGroupSize)
	// All bits must be aligned in groups by 3, so to keep this alignment need to add special bit,
	// after the most left block and by doing so number will look like this:
	//
	// 0b1<owner-group><shared-group><other-group>
	//   ↑
	//   alignmentBit
	//
	// E.g: 0b1_111_011_001
	return FilePermissions(((alignmentBit | (owner << offset) | shared) << offset) | other)
}

var emptyFilePermissions = NewFilePermissions(0, 0, 0)

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
	otherGroupROffset int = FilePermissionGroupSize * iota
	sharedGroupROffset
	ownerGroupROffset
)

func (p FilePermissions) getGroupWithROffset(rightOffset int) FilePermissionGroup {
	return FilePermissionGroup(p >> rightOffset & filePermissionsGroupBits)
}

func (p FilePermissions) GetOwnerPermissions() FilePermissionGroup {
	return p.getGroupWithROffset(ownerGroupROffset)
}

func (p FilePermissions) GetSharedPermissions() FilePermissionGroup {
	return p.getGroupWithROffset(sharedGroupROffset)
}

func (p FilePermissions) GetOtherPermissions() FilePermissionGroup {
	return p.getGroupWithROffset(otherGroupROffset)
}

func (p FilePermissions) String() string {
	str := make([]string, FilePermissionGroupsAmount)

	for i := range FilePermissionGroupsAmount {
		str[FilePermissionGroupsAmount-i-1] = p.getGroupWithROffset(FilePermissionGroupSize * i).String()
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

	// TODO probably this should be optimised... at least to be more readable

	read := rune(group[0])
	update := rune(group[1])
	del := rune(group[2])

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
			rawGroups[rawGroupIdx] = permissions[offset : FilePermissionGroupSize*(rawGroupIdx+1)]
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

var ChecksumHasher hash.Hash = sha256.New()

type UpdatableFileMetadata struct {
	Bucket string
	Path   string

	Encoding string

	Owner       string
	Permissions FilePermissions

	Description string
	Categories  []string
	Tags        []string

	Status FileStatus
}

type GeneratedFileMetadata struct {
	ID   		 string // UUID string
	OriginalName string

	MIMEType     string
	Size         int64
	Checksum     string

	UploadedBy string

	UploadedAt time.Time
	UpdatedAt  time.Time
	CreatedAt  time.Time
	AccessedAt time.Time
}

type FileMetadata struct {
	UpdatableFileMetadata
	GeneratedFileMetadata
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
