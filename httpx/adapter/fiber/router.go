package fiber

import (
	"github.com/cobanhub/lib/httpx"
	"github.com/gofiber/fiber/v3"
)

type router struct {
	app *fiber.App
}

var _ httpx.Router = (*router)(nil)

func New() httpx.Router {
	return &router{
		app: fiber.New(),
	}
}
