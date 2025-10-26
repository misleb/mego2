//go:build js && wasm

package client

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/misleb/mego2/shared"
)

// replacePathParams replaces path parameters in URL with values from request
func replacePathParams(url string, request interface{}) string {
	requestValue := reflect.ValueOf(request)
	requestType := reflect.TypeOf(request)

	for i := 0; i < requestValue.NumField(); i++ {
		field := requestType.Field(i)
		value := requestValue.Field(i)

		// Check for URI tag
		if uriTag := field.Tag.Get("uri"); uriTag != "" {
			paramValue := fmt.Sprintf("%v", value.Interface())
			url = strings.Replace(url, ":"+uriTag, paramValue, 1)
		}
	}

	return url
}

// Couldn't be made part of APIClient because methods can't take type parameters in Go 1.25.1
func CallEndpointTyped[T any](c *APIClient, endpoint shared.Endpoint, request interface{}) (*T, error) {
	response, err := c.CallEndpoint(endpoint, request)
	if err != nil {
		return nil, err
	}

	typedResp, ok := response.(*T)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return typedResp, nil
}
