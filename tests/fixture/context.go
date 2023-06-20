package fixture

import (
	"context"
	"time"
)

// UserID is a fixture of schema entity for testing
var UserID = map[string]int{
	"admin":        1,
	"android_user": 215,
	"ios_user":     216,
}

// CtxEnded creates dummy context with cancelled state.
func CtxEnded() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond))
	defer cancel()
	return ctx
}
