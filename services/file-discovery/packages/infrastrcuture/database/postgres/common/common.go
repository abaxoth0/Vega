package dbcommon

import (
	"errors"
	"fmt"
	"vega_file_discovery/packages/entity"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
)

// Converts files status from it's plain representation, to the format in which it's stored in DB.
// Panics if files status is invalid.
//
// * ActiveFileStatus - 'A'
//
// * PendingFilestatus - 'P'
//
// * ArchivedFileStatus - 'R'
//
func ConvertFileStatus(status entity.FileStatus) string {
	switch status {
	case entity.ActiveFileStatus:
		return "A"
	case entity.PendingFilestatus:
		return "P"
	case entity.ArchivedFileStatus:
		return "R"
	default:
		dblog.Logger.Panic(
			"Failed to convert files status for DB",
			fmt.Sprintf("Unknown files status '%v'", status),
			nil,
		)
	}
	return ""
}

var ErrInvalidRawFileStatus = errors.New("Invalid raw files status")

// Parses raw files status to it's plain representation.
// Panics if files status is invalid.
//
// * ActiveFileStatus - 'A'
//
// * PendingFilestatus - 'P'
//
// * ArchivedFileStatus - 'R'
//
func ParseFileStatus(rawFileStatus string) (entity.FileStatus, error) {
	switch rawFileStatus {
	case "A":
		return entity.ActiveFileStatus, nil
	case "P":
		return entity.PendingFilestatus, nil
	case "R":
		return entity.ArchivedFileStatus, nil
	default:
		return 0, ErrInvalidRawFileStatus
	}
}
