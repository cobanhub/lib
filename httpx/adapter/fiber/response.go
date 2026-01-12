package fiber

import "github.com/gofiber/fiber/v3"

type fiberResponse struct {
	ctx fiber.Ctx
}

func (r *fiberResponse) Status(code int) {
	r.ctx.Status(code)
}

func (r *fiberResponse) JSON(code int, body any) error {
	return r.ctx.Status(code).JSON(body)
}
