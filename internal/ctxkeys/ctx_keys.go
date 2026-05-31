package ctxkeys

type contextKey struct{}

var (
	CtxUserIDKey    contextKey = struct{}{}
	CtxRealIPKey    contextKey = struct{}{}
	CtxRequestIDKey contextKey = struct{}{}
)
