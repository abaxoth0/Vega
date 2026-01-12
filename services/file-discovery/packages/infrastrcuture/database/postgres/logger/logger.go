package log

import "github.com/abaxoth0/Vega/libs/go/packages/logger"

var (
	DB        = logger.NewSource("DATABASE", logger.Default)
	Migration = logger.NewSource("MIGRATION", logger.Default)
)

