// Command Query Responsibility Segregation
// (https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)
package cqrs

import (
	"context"
	"time"
)

const DefaultCommandQueryTimeout time.Duration = time.Second * 3

// Base of CQRS Commands and Queries structs.
// This structs are used to to encapsulate the intent of an operation.
//
// Commands are responsible for write operations.
// Queries are responsible for read operations.
type CommandQuery struct {
	Context        context.Context
	ContextTimeout time.Duration
}

// Verifies that both context and it's timeout are exist and valid.
// CommandQuery considered invalid if: Context is nil OR timeout is <= 0
func (c CommandQuery) IsInit() bool {
	return c.Context != nil && c.ContextTimeout > 0
}

// Creates new background context for Context field if it's nil.
//
// Also sets ContextTimeout to DefaultCommandQueryTimeout if it's <= 0.
func InitDefaultCommandQuery(c *CommandQuery) {
	if c.Context == nil {
		c.Context = context.Background()
	}
	if c.ContextTimeout <= 0 {
		c.ContextTimeout = DefaultCommandQueryTimeout
	}
}

type IdTargetedCommandQuery struct {
	ID string

	CommandQuery
}
