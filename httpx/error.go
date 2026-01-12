// httpx/errors.go
package httpx

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func HandleError(ctx Context, err error) {
	if err == nil {
		return
	}

	if he, ok := err.(*HTTPError); ok {
		_ = ctx.JSON(he.Code, map[string]string{
			"error": he.Message,
		})
		return
	}

	_ = ctx.JSON(500, map[string]string{
		"error": "internal server error",
	})
}
