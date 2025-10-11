package endpoint

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/misleb/mego2/shared"
)

// Generic endpoint registration function
func RegisterEndpoint(router *gin.Engine, endpoint shared.Endpoint, handler interface{}) {
	var handlerFuncs []gin.HandlerFunc

	if endpoint.AuthRequired {
		handlerFuncs = append(handlerFuncs, authRequired())
	}

	handlerFuncs = append(handlerFuncs, createHandler(endpoint, handler))

	switch endpoint.Method {
	case http.MethodGet:
		router.GET(endpoint.Path, handlerFuncs...)
	case http.MethodPost:
		router.POST(endpoint.Path, handlerFuncs...)
	case http.MethodPut:
		router.PUT(endpoint.Path, handlerFuncs...)
	case http.MethodDelete:
		router.DELETE(endpoint.Path, handlerFuncs...)
	}
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")

		if token != "currToken" {
			c.JSON(401, shared.LoginResponse{Error: "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func createHandler(endpoint shared.Endpoint, handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new instance of the request type
		requestType := reflect.TypeOf(endpoint.RequestType)
		requestValue := reflect.New(requestType).Interface()

		// Bind URI parameters first
		if err := c.ShouldBindUri(requestValue); err != nil {
			sendErrorResponse(c, endpoint, err)
			return
		}

		// For POST/PUT requests, also bind JSON body
		if endpoint.Method == http.MethodPost || endpoint.Method == http.MethodPut {
			if err := c.ShouldBindJSON(requestValue); err != nil {
				sendErrorResponse(c, endpoint, err)
				return
			}
		}

		// Call the handler with the bound request
		handlerValue := reflect.ValueOf(handler)
		handlerValue.Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(requestValue).Elem(),
		})
	}
}

func sendErrorResponse(c *gin.Context, endpoint shared.Endpoint, err error) {
	responseType := reflect.TypeOf(endpoint.ResponseType)
	errorResponse := reflect.New(responseType).Interface()

	// Set error field if it exists
	if errorField := reflect.ValueOf(errorResponse).Elem().FieldByName("Error"); errorField.IsValid() && errorField.CanSet() {
		errorField.SetString(err.Error())
	}

	c.JSON(400, errorResponse)
}
