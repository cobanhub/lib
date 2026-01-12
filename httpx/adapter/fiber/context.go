package fiber

import (
	"github.com/gofiber/fiber/v3"
)

type fiberCtx struct {
	ctx fiber.Ctx
}

func (f fiberCtx) Param(key string) string {
	return f.ctx.Params(key)
}

func (f *fiberCtx) Query(key string) string {
	return f.ctx.Query(key)
}

func (f *fiberCtx) Queries() map[string]string {
	return f.ctx.Queries()
}

func (f *fiberCtx) Bind(v any) error {
	return f.ctx.Bind().Body(v)
}

func (f *fiberCtx) JSON(code int, body any) error {
	return f.ctx.Status(code).JSON(body)
}
