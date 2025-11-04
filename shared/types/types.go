package types

import "net/http"

type Password string

type RequestAugment func(*http.Request) error

type GoogleTokenExchangeRequest struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GoogleTokenExchangeResponse struct {
	IDToken          string `json:"id_token"`
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
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

type User struct {
	Name  string
	Token string
	Email string
}
