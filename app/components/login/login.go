//go:build js && wasm

package login

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/basic/button"
	"github.com/gofred-io/gofred/breakpoint"
	"github.com/gofred-io/gofred/console"
	"github.com/gofred-io/gofred/foundation/column"
	"github.com/gofred-io/gofred/foundation/container"
	"github.com/gofred-io/gofred/foundation/spacer"
	"github.com/gofred-io/gofred/foundation/text"
	"github.com/gofred-io/gofred/options/spacing"
	"github.com/misleb/mego2/app/client"
	"github.com/misleb/mego2/app/components/basic/input"
	"github.com/misleb/mego2/app/store"
	"github.com/misleb/mego2/shared/api_client"
	"github.com/misleb/mego2/shared/types"
)

var (
	userInput *input.Input
	passInput *input.Input
)

func Get() application.BaseWidget {
	if userInput == nil {
		userInput = input.New("", input.Placeholder("Username"))
		passInput = input.New("", input.Placeholder("Password"), input.Type(input.InputTypePassword))
	}
	return container.New(
		column.New(
			[]application.BaseWidget{
				text.New("Who Are You?"),
				spacer.New(spacer.Height(20)),
				userInput.BaseWidget,
				spacer.New(spacer.Height(20)),
				passInput.BaseWidget,
				spacer.New(spacer.Height(20)),
				button.New(text.New("Login"), button.OnClick(doLogin)),
				spacer.New(spacer.Height(10)),
				button.New(text.New("Sign in with Google"), button.OnClick(doGoogleLogin)),
			},
		),
		container.Padding(breakpoint.All(spacing.All(32))),
	)
}

func doLogin(this application.BaseWidget, e application.Event) {
	c := api_client.GetInstance()
	result, err := c.Login(userInput.GetValue(), passInput.GetValue())
	if err != nil {
		console.Error(err.Error())
		return
	}
	store.SetUser(result.User)
}

func doGoogleLogin(this application.BaseWidget, e application.Event) {
	// Get auth code from JavaScript
	code, err := client.InitiateGoogleAuth()
	if err != nil {
		console.Error("Failed to initiate Google auth: " + err.Error())
		return
	}

	// Send code to backend
	apiClient := api_client.GetInstance()
	request := types.GoogleAuthRequest{Code: code}
	response, err := api_client.CallEndpointTyped[types.LoginResponse](
		apiClient,
		types.GoogleAuthEndpoint,
		request,
		api_client.NoOpRequestAugment,
	)

	if err != nil {
		console.Error("Failed to authenticate with Google: " + err.Error())
		return
	}

	// Store user
	store.SetUser(response.User)
}
