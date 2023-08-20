package response_handling

import (
	"io"
	"net/http"
)

func GetBody(ctx *HandlingContext, resp *http.Response) (string, error) {
	bodyP := &ctx.ResponseBody
	if *bodyP == "" {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		*bodyP = string(bodyBytes)
	}
	return *bodyP, nil
}

func DefaultRegistrationHandler(agentId *string) ResponseHandler {
	return defaultRequestHandler([]Mapping{{"debuggee.id", agentId}})
}

func DefaultAuthenticationHandler(accessToken *string) ResponseHandler {
	return defaultRequestHandler([]Mapping{{"id_token", accessToken}})
}

func defaultRequestHandler(mappings []Mapping) ResponseHandler {
	contextP := &HandlingContext{}
	return &CompositeHandler{
		HandlersChain: []ResponseHandler{
			&StatusHandler{contextP},
			&ParsingHandler{
				Mappings: mappings,
				Context:  contextP,
			},
		},
	}
}
