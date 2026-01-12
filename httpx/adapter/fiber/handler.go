package fiber

import (
	"github.com/gofiber/fiber/v3"

	"github.com/cobanhub/lib/httpx"
)

func (r *router) handle(
	method string,
	path string,
	h httpx.Handler,
) {
	r.app.Add([]string{method}, path, func(c fiber.Ctx) error {
		req := &fiberRequest{ctx: c}
		res := &fiberResponse{ctx: c}

		if err := h(c.Context(), req, res); err != nil {
			handleError(res, err)
		}
		return nil
	})
}
