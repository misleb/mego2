package shared

import "net/http"

var IncEndpoint = Endpoint{
	Path:         "/inc/:value",
	Method:       http.MethodGet,
	RequestType:  IntRequest{},
	ResponseType: IntResponse{},
	AuthRequired: true,
}

var LoginEndpoint = Endpoint{
	Path:         "/login",
	Method:       http.MethodPost,
	RequestType:  LoginRequest{},
	ResponseType: LoginResponse{},
	AuthRequired: false,
}
