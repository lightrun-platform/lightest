package response_handling

import (
	"net/http"
)

type ResponseHandler interface {
	Handle(res *http.Response) error
}

type HandlingContext struct {
	ResponseBody string
}

func (c HandlingContext) GetContext() *HandlingContext {
	return &c
}

type HandlingError struct {
	Message string
}

func (e HandlingError) Error() string {
	return e.Message
}
