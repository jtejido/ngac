package context

import (
    "context"
    "sync"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type ngacContextKey struct {
    name string
}

var (
    // ContextKeySessionID is a context key for use with Contexts in this package.
    // The associated value will be of type int.
    ContextKeySession = &ngacContextKey{"session-id"}

    // ContextKeyProcessID is a context key for use with Contexts in this package.
    // The associated value will be of type int.
    ContextKeyProcess = &ngacContextKey{"process-id"}
)

// Context is a PDP specific context interface. It's used in
// authentication handlers and callbacks, and its underlying context.Context is
// exposed on Session in the session Handler. A connection-scoped lock is also
// embedded in the context to make it easier to limit operations per-connection.
type Context interface {
    context.Context
    sync.Locker

    // UserID returns the user's ID.
    User() string

    // ProcessID returns the process ID.
    Process() string

    // SetValue allows you to easily write new values into the underlying context.
    SetValue(key, value interface{})
}

type userContext struct {
    context.Context
    *sync.Mutex
}

func newUserContext() (*userContext, context.CancelFunc) {
    innerCtx, cancel := context.WithCancel(context.Background())
    ctx := &userContext{innerCtx, &sync.Mutex{}}
    return ctx, cancel
}

func NewUserContext(user string) (*userContext, context.CancelFunc) {
    ctx, cancel := newUserContext()
    ctx.SetValue(ContextKeySession, user)
    ctx.SetValue(ContextKeyProcess, "")
    return ctx, cancel
}

func NewUserContextWithProcess(user, process string) (*userContext, context.CancelFunc) {
    ctx, cancel := newUserContext()
    ctx.SetValue(ContextKeySession, user)
    ctx.SetValue(ContextKeyProcess, process)
    return ctx, cancel
}

func (ctx *userContext) SetValue(key, value interface{}) {
    ctx.Context = context.WithValue(ctx.Context, key, value)
}

func (ctx *userContext) User() string {
    return ctx.Value(ContextKeySession).(string)
}

func (ctx *userContext) Process() string {
    return ctx.Value(ContextKeyProcess).(string)
}
