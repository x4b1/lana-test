package log

import "context"

// Logger provides loggin functions
type Logger interface {
	Error(ctx context.Context, err error)
	Fatal(ctx context.Context, err error)
	Info(ctx context.Context, msg string, args ...interface{})
}
