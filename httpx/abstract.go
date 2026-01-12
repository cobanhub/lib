package httpx

type Request interface {
	Param(string) string
	Query(string) string
	Header(string) string
	Bind(any) error
}

type Response interface {
	JSON(status int, body any) error
	Status(code int)
}
