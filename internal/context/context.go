package context

// contextKey is an unexported type to prevent collisions.
type contextKey string

const (
	UserIDContextKey contextKey = "user_id"
)
