package fiber

import (
	"github.com/cobanhub/lib/httpx"
	"github.com/gofiber/fiber/v3"
)

func (r *router) GET(path string, h httpx.Handler) {
	r.handle(fiber.MethodGet, path, h)
}

func (r *router) POST(path string, h httpx.Handler) {
	r.handle(fiber.MethodPost, path, h)
}

func (r *router) PUT(path string, h httpx.Handler) {
	r.handle(fiber.MethodPut, path, h)
}

func (r *router) DELETE(path string, h httpx.Handler) {
	r.handle(fiber.MethodDelete, path, h)
}

func (r *router) Start(addr string) error {
	return r.app.Listen(addr)
}
