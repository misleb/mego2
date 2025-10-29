package shared

type Password string

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
