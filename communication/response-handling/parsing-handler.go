package response_handling

import (
	"github.com/tidwall/gjson"
	"net/http"
)

// ParsingHandler A response handler that parses JSON paths into given fields
type ParsingHandler struct {
	Mappings []Mapping
	Context  *HandlingContext
}

type Mapping struct {
	Path   string
	Target *string
}

func (handler *ParsingHandler) Handle(resp *http.Response) error {
	body := handler.Context.ResponseBody
	if body == "" {
		return HandlingError{Message: "Getting response body failed."}
	}
	for _, m := range handler.Mappings {
		*m.Target = gjson.Get(body, m.Path).String()
	}
	return nil
}
