//go:build js && wasm

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/misleb/mego2/app/store"
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
		url = replacePathParams(url, request)
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
		if user := store.GetUser(); user != nil {
			req.Header.Add("X-Auth-Token", user.Token)
		} else {
			return nil, fmt.Errorf("No user to get token from for athenticated request. Aborting.")
		}

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

// Convenience methods for specific endpoints
func (c *APIClient) Login(user string, pass string) (*shared.LoginResponse, error) {
	request := shared.LoginRequest{Username: user, Password: pass}
	return CallEndpointTyped[shared.LoginResponse](c, shared.LoginEndpoint, request)
}

// Increment calls the increment endpoint
func (c *APIClient) Increment(value int) (*shared.IntResponse, error) {
	request := shared.IntRequest{Value: value}
	return CallEndpointTyped[shared.IntResponse](c, shared.IncEndpoint, request)
}
