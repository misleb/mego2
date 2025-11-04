//go:build js && wasm

package client

import (
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/misleb/mego2/app/store"
)

var AddAuthHeader = func(req *http.Request) error {
	if user := store.GetUser(); user != nil {
		req.Header.Add("X-Auth-Token", user.Token)
		return nil
	} else {
		return fmt.Errorf("no user to get token from")
	}
}

// InitiateGoogleAuth starts the Google OAuth flow via JavaScript
func InitiateGoogleAuth() (string, error) {
	// Call JavaScript function
	jsFunc := js.Global().Get("initiateGoogleAuth")
	if jsFunc.IsUndefined() {
		return "", fmt.Errorf("initiateGoogleAuth function not found")
	}

	// Create a promise channel
	promiseChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Call the JS function which returns a Promise
	promise := jsFunc.Invoke()

	// Set up promise handlers
	then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			code := args[0].String()
			promiseChan <- code
		} else {
			errorChan <- fmt.Errorf("no code returned")
		}
		return nil
	})
	defer then.Release()

	catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			errMsg := args[0].Get("message").String()
			errorChan <- fmt.Errorf("JavaScript error: %s", errMsg)
		} else {
			errorChan <- fmt.Errorf("unknown JavaScript error")
		}
		return nil
	})
	defer catch.Release()

	promise.Call("then", then).Call("catch", catch)

	// Wait for result
	select {
	case code := <-promiseChan:
		return code, nil
	case err := <-errorChan:
		return "", err
	}
}
