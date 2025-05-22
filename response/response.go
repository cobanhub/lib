package response

import (
	"github.com/gofiber/fiber/v3"
)

type JSONResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{Code: fiber.StatusOK, Data: nil, Message: ""}
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetCode(code int) *JSONResponse {
	r.Code = code
	return r
}
