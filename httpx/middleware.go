// middleware.go
package httpx

import "context"

type Handler func(ctx context.Context, req Request, res Response) error
