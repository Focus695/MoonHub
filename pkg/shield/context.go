package shield

import "context"

type ctxKeyApprovedTool struct{}

// ContextWithApprovedToolExecution marks ctx so tools can skip in-Execute shield
// re-checks after the agent loop has obtained user approval for the same action.
func ContextWithApprovedToolExecution(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyApprovedTool{}, true)
}

// ApprovedToolExecution reports whether the current tool execution was already
// approved by the user at the agent loop (see ContextWithApprovedToolExecution).
func ApprovedToolExecution(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyApprovedTool{}).(bool)
	return v
}
