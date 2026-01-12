// router.go
package httpx

import "context"

type Handler func(ctx context.Context, req Request, res Response) error

type Router interface {
	GET(path string, h Handler)
	POST(path string, h Handler)
	PUT(path string, h Handler)
	DELETE(path string, h Handler)

	Start(addr string) error
}
