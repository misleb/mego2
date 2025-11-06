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

type GoogleAuthRequest struct {
	Code string `json:"code"`
}

type GoogleAuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type IntResponse struct {
	Result int    `json:"result"`
	Error  string `json:"error"`
}

type IntRequest struct {
	Value int `uri:"value" json:"value" form:"value"`
}

type Endpoint struct {
	Path         string
	Method       string
	RequestType  interface{}
	ResponseType interface{}
	AuthRequired bool
}
