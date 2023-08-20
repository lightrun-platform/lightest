package requests

import "strings"

var agentRegistration string = "debuggees/company/{companyId}/register/"
var getBreakpoints string = "debuggees/company/{companyId}/{agentId}/breakpoints?successOnTimeout=false"
var authenticate string = "api/authenticate/"
var websocketConnection string = "socket/events?access_token={accessToken}&company={companyId}"

func ProduceAgentRegistrationURL(serverUrl string, companyId string) string {
	return serverUrl + strings.Replace(agentRegistration, "{companyId}", companyId, -1)
}

func ProduceGetBreakpointsURL(serverUrl string, companyId string, agentId string) string {
	r := strings.NewReplacer("{companyId}", companyId, "{agentId}", agentId)
	return serverUrl + r.Replace(getBreakpoints)
}

func ProduceAuthenticateURL(serverUrl string) string {
	return serverUrl + authenticate
}

func ProduceWebsocketConnectionURL(serverUrl string, accessToken string, companyId string) string {
	r := strings.NewReplacer("{companyId}", companyId, "{accessToken}", accessToken)
	return strings.Replace(serverUrl, "https", "wss", 1) + r.Replace(websocketConnection)
}
