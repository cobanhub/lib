package fiber

import (
	"github.com/cobanhub/lib/httpx"
)

func handleError(res httpx.Response, err error) {
	if err == nil {
		return
	}

	if he, ok := err.(*httpx.HTTPError); ok {
		_ = res.JSON(he.Code, map[string]string{
			"error": he.Message,
		})
		return
	}

	_ = res.JSON(500, map[string]string{
		"error": "internal server error",
	})
}
