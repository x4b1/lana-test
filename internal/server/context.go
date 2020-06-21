package server

import "context"

var (
	contextKeyEndpoint = contextKey("endpoint")
)

type contextKey string

func (c contextKey) String() string {
	return "server" + string(c)
}

//Endpoint adds endpoint path to context
func Endpoint(ctx context.Context) (string, bool) {
	endpoint, ok := ctx.Value(contextKeyEndpoint).(string)
	return endpoint, ok
}
