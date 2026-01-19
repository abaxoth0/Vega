package filemetadatatable

import (
	"fmt"
	"vega_file_discovery/packages/entity"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
)

type Manager struct {
	//
}

// Converts files status from it's plain representation, to the format in which it's stored in DB.
// Panics if files status is invalid.
//
// * ActiveFileStatus - 'A'
//
// * PendingFilestatus - 'P'
//
// * ArchivedFileStatus - 'R'
//
func convertFileStatus(status entity.FileStatus) string {
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
