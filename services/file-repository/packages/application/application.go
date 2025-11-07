package application

import (
	"context"
	"time"
)

type CommandQuery struct {
	Context        context.Context
	ContextTimeout time.Duration
}

// Verifies that both context and it's timeout are exist and valid.
// CommandQuery considered invalid if: Context is nil OR timeout is <= 0
func (c CommandQuery) IsInit() bool {
	return c.Context != nil && c.ContextTimeout > 0
}

// Creates new background context for Context fild and sets timeout to default (now it's 3 seconds)
func InitDefaultCommandQuery(c *CommandQuery) {
	if c.Context == nil {
		c.Context = context.Background()
	}
	if c.ContextTimeout <= 0 {
		// TODO move right part into the constant
		c.ContextTimeout = time.Second * 3
	}
}
