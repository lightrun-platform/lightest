package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	. "prerequisites-tester/communication/response-handling"
	. "prerequisites-tester/config"
)

type RegistrationData struct {
	AgentConfig
	ServerUrl string
	CompanyId string
}

func RegisterAgent(data RegistrationData, handler ResponseHandler) error {
	requestUrl := ProduceAgentRegistrationURL(data.ServerUrl, data.CompanyId)
	payloadBytes, err := buildRegistrationPayloadBytes(data)
	if err != nil {
		return err
	}

	return makeRequest(
		constructRequestProducer("POST", requestUrl, payloadBytes),
		handler,
	)
}

func GetBreakpoints(agentId string, handler ResponseHandler) error {
	requestUrl := ProduceGetBreakpointsURL(GetConfig().ServerUrl, GetConfig().CompanyId, agentId)

	return makeRequest(
		constructRequestProducer("GET", requestUrl, nil),
		handler,
	)
}

func Authenticate(userEmail string, userPassword string, handler ResponseHandler) error {
	requestUrl := ProduceAuthenticateURL(GetConfig().ServerUrl)
	payloadBytes, err := buildAuthenticationPayloadBytes(userEmail, userPassword)
	if err != nil {
		return err
	}

	return makeRequest(
		constructRequestProducer("POST", requestUrl, payloadBytes),
		handler,
	)
}

func constructRequestProducer(
	requestType string, requestUrl string, requestPayload []byte) func() (*http.Request, error) {
	var payloadBuffer io.Reader = nil
	if requestPayload != nil {
		payloadBuffer = bytes.NewBuffer(requestPayload)
	}
	return func() (*http.Request, error) {
		return http.NewRequest(requestType, requestUrl, payloadBuffer)
	}
}

func makeRequest(requestProducer func() (*http.Request, error), handler ResponseHandler) error {
	// Produce request entity.
	req, err := requestProducer()
	if err != nil {
		return err
	}

	// Attach headers to the request entity.
	attachHeaders(req)

	// Create a client and use it to make the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Handle the response.
	return handler.Handle(resp)
}

func attachHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", GetConfig().ApiKey)
}

func buildRegistrationPayloadBytes(data RegistrationData) ([]byte, error) {
	debuggeePayload := map[string]interface{}{
		"agentVersion":              data.Version,
		"host":                      data.HostName,
		"id":                        "",
		"pid":                       os.Getpid(),
		"uniquifier":                "",
		"isDisabled":                false,
		"cpuCoreCount":              0,
		"systemTotalMemoryBytes":    0,
		"runtimeEnvironment":        "Java",
		"agentOS":                   "Ubuntu",
		"osVersion":                 "",
		"linuxDistroName":           "",
		"runtimeEnvironmentVersion": "",
		"runtimeEnvironmentInfo":    "",
		"procArch":                  "",
		"isInDocker":                false,
		"isInK8s":                   false,
	}
	payload := map[string]interface{}{
		"debuggee": debuggeePayload,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, HandlingError{Message: fmt.Sprintf(
			"Failed to produce agent registration request's JSON body:", err,
		)}
	}
	return payloadBytes, nil
}

func buildAuthenticationPayloadBytes(userEmail string, userPassword string) ([]byte, error) {
	payload := map[string]interface{}{
		"email":      userEmail,
		"password":   userPassword,
		"rememberMe": true,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, HandlingError{Message: fmt.Sprintf(
			"Failed to produce client authentication request's JSON body:", err,
		)}
	}
	return payloadBytes, nil
}
