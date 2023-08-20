package response_handling

import (
	"net/http"
)

type CompositeHandler struct {
	HandlersChain []ResponseHandler
}

func (handler *CompositeHandler) Handle(resp *http.Response) error {
	for _, nextHandler := range handler.HandlersChain {
		err := nextHandler.Handle(resp)
		if err != nil {
			return err
		}
	}
	return nil
}
