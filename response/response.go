package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceID string      `json:"trace_id,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewResponseSuccess(meta *Pagination, data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "OK",
		Meta:    meta,
		Data:    data,
	}
}

func NewResponseCreated(data interface{}) *Response {
	return &Response{
		Code:    http.StatusCreated,
		Message: "CREATED",
		Data:    data,
	}
}

func NewResponseInvalidRequest(data interface{}) *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Message: "INVALID_REQUEST",
		Data:    data,
	}
}

func NewResponseBadRequest(message string) *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Message: "BAD_REQUEST",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseNotFound(message string) *Response {
	return &Response{
		Code:    http.StatusNotFound,
		Message: "NOT_FOUND",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseUnauthorized(message string) *Response {
	return &Response{
		Code:    http.StatusUnauthorized,
		Message: "UNAUTHORIZED",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseForbidden(message string) *Response {
	return &Response{
		Code:    http.StatusForbidden,
		Message: "FORBIDDEN",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseInternalServerError(message string) *Response {
	return &Response{
		Code:    http.StatusInternalServerError,
		Message: "INTERNAL_SERVER_ERROR",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseConflict(message string) *Response {
	return &Response{
		Code:    http.StatusConflict,
		Message: "CONFLICT",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseConflictWithData(message string, data interface{}) *Response {
	return &Response{
		Code:    http.StatusConflict,
		Message: message,
		Data:    data,
	}
}

func NewResponseValidationError(errors map[string]string) *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Message: "VALIDATION_ERROR",
		Data: map[string]interface{}{
			"errors": errors,
		},
	}
}

func NewResponseUnprocessableEntity(message string) *Response {
	return &Response{
		Code:    http.StatusUnprocessableEntity,
		Message: "UNPROCESSABLE_ENTITY",
		Data: map[string]string{
			"message": message,
		},
	}
}

func NewResponseWithTraceID(code int, message string, data interface{}, traceID string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		TraceID: traceID,
		Data:    data,
	}
}

func (r *Response) Encode() []byte {
	res, err := json.Marshal(r)
	if err != nil {
		fmt.Println("failed to encode message")
		return nil
	}
	return res
}
