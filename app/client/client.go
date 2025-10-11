package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/misleb/mego2/shared"
)

// APIClient provides methods to call API endpoints
type APIClient struct {
	baseURL string
	client  *http.Client
}

var (
	instance *APIClient
	once     sync.Once
)

// GetInstance returns the singleton instance of APIClient
func GetInstance() *APIClient {
	once.Do(func() {
		instance = &APIClient{
			baseURL: "", // Default to relative paths
			client:  &http.Client{},
		}
	})
	return instance
}

// SetBaseURL allows you to configure the base URL for the singleton
func (c *APIClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// CallEndpoint makes a generic call to any endpoint defined in shared
func (c *APIClient) CallEndpoint(endpoint shared.Endpoint, request interface{}) (interface{}, error) {
	// Build the URL with path parameters
	url := c.baseURL + endpoint.Path

	// Replace path parameters if request has URI tags
	if request != nil {
		url = c.replacePathParams(url, request)
	}

	// Prepare the request body for POST/PUT
	var body []byte
	var err error
	if (endpoint.Method == http.MethodPost || endpoint.Method == http.MethodPut) && request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(endpoint.Method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if endpoint.AuthRequired {
		req.Header.Add("X-Auth-Token", "currToken")
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Make the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Create response instance
	responseType := reflect.TypeOf(endpoint.ResponseType)
	responseValue := reflect.New(responseType).Interface()

	// Decode response
	if err := json.NewDecoder(resp.Body).Decode(responseValue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for error in response
	if resp.StatusCode >= 400 {
		// Try to extract error message if available
		if errorField := reflect.ValueOf(responseValue).Elem().FieldByName("Error"); errorField.IsValid() {
			if errorMsg := errorField.String(); errorMsg != "" {
				return nil, fmt.Errorf("API error: %s", errorMsg)
			}
		}
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return responseValue, nil
}

// replacePathParams replaces path parameters in URL with values from request
func (c *APIClient) replacePathParams(url string, request interface{}) string {
	requestValue := reflect.ValueOf(request)
	requestType := reflect.TypeOf(request)

	for i := 0; i < requestValue.NumField(); i++ {
		field := requestType.Field(i)
		value := requestValue.Field(i)

		// Check for URI tag
		if uriTag := field.Tag.Get("uri"); uriTag != "" {
			paramValue := fmt.Sprintf("%v", value.Interface())
			url = replaceParam(url, ":"+uriTag, paramValue)
		}
	}

	return url
}

// replaceParam replaces a parameter placeholder in URL with actual value
func replaceParam(url, param, value string) string {
	// Simple replacement - you might want to use url.PathEscape for proper encoding
	return strings.Replace(url, param, value, 1)
}

// Convenience methods for specific endpoints

// Increment calls the increment endpoint
func (c *APIClient) Increment(value int) (*shared.IntResponse, error) {
	request := shared.IntRequest{Value: value}
	response, err := c.CallEndpoint(shared.IncEndpoint, request)
	if err != nil {
		return nil, err
	}

	// Type assert to the expected response type
	intResp, ok := response.(*shared.IntResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return intResp, nil
}
