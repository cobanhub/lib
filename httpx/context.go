package httpx

import "context"

type Context interface {
	Context() context.Context

	Param(key string) string
	Query(key string) string
	Header(key string) string
	Queries() map[string]string

	Bind(any) error
	JSON(status int, body any) error

	Set(key string, val any)
	Get(key string) (any, bool)
}
