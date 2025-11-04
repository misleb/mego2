package types

import (
	"net/http"
)

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

var GoogleAuthEndpoint = Endpoint{
	Path:         "/auth/google",
	Method:       http.MethodPost,
	RequestType:  GoogleAuthRequest{},
	ResponseType: GoogleAuthResponse{},
	AuthRequired: false,
}

var GoogleTokenExchangeEndpoint = Endpoint{
	Path:         "https://oauth2.googleapis.com/token",
	Method:       http.MethodPost,
	RequestType:  GoogleTokenExchangeRequest{},
	ResponseType: GoogleTokenExchangeResponse{},
	AuthRequired: false,
}
