package fiber

import "github.com/gofiber/fiber/v3"

type fiberRequest struct {
	ctx fiber.Ctx
}

func (r *fiberRequest) Param(key string) string {
	return r.ctx.Params(key)
}

func (r *fiberRequest) Query(key string) string {
	return r.ctx.Query(key)
}

func (r *fiberRequest) Header(key string) string {
	return r.ctx.Get(key)
}

func (r *fiberRequest) Bind(v any) error {
	return r.ctx.Bind().Body(v)
}
