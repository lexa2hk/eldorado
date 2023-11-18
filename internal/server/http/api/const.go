package api

const (
	RequestIDHeader     = "X-Request-ID"
	AuthorizationHeader = "Authorization"
)

type key int

const (
	UserIDKey    key = 0
	RequestIDKey key = 1
	TokenKey     key = 2
)
