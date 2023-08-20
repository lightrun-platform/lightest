package response_handling

import (
	"io"
	"net/http"
)

type StatusHandler struct {
	Context *HandlingContext
}

func (handler *StatusHandler) Handle(resp *http.Response) error {
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return HandlingError{Message: "Request failed with status:" + resp.Status}
	}

	// Saving the response body.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body := string(bodyBytes)
	handler.Context.ResponseBody = string(body)

	return nil
}
