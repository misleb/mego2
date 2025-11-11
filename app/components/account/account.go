//go:build js && wasm

package account

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/basic/button"
	"github.com/gofred-io/gofred/breakpoint"
	"github.com/gofred-io/gofred/console"
	"github.com/gofred-io/gofred/foundation/column"
	"github.com/gofred-io/gofred/foundation/container"
	"github.com/gofred-io/gofred/foundation/spacer"
	"github.com/gofred-io/gofred/foundation/text"
	"github.com/gofred-io/gofred/listenable"
	"github.com/gofred-io/gofred/options/spacing"
	"github.com/misleb/mego2/app/client"
	"github.com/misleb/mego2/app/components/basic/input"
	"github.com/misleb/mego2/app/store"
	"github.com/misleb/mego2/shared/api_client"
	"github.com/misleb/mego2/shared/types"
)

var (
	newPassInput *input.Input
	oldPassInput *input.Input
)

func Get() application.BaseWidget {
	if newPassInput == nil {
		newPassInput = input.New("", input.Placeholder("New Password"), input.Type(input.InputTypePassword))
		oldPassInput = input.New("", input.Placeholder("Old Password"), input.Type(input.InputTypePassword))
	}
	return container.New(
		column.New(
			[]application.BaseWidget{
				text.New("Account"),
				spacer.New(spacer.Height(20)),
				listenable.Builder(store.AppStoreListenable(), func() application.BaseWidget {
					user := store.GetUser()
					if user.IsNewExternal {
						return text.New("You are a new external user. Please set a password to continue.")
					} else {
						return oldPassInput.BaseWidget
					}
				}),
				spacer.New(spacer.Height(20)),
				newPassInput.BaseWidget,
				spacer.New(spacer.Height(20)),
				button.New(text.New("Save"), button.OnClick(doAccountUpdate)),
			},
		),
		container.Padding(breakpoint.All(spacing.All(32))),
	)
}

func doAccountUpdate(this application.BaseWidget, e application.Event) {
	apiClient := api_client.GetInstance()
	request := types.UpdateSelfRequest{
		OldPassword: oldPassInput.GetValue(),
		NewPassword: newPassInput.GetValue(),
	}
	result, err := api_client.CallEndpointTyped[types.LoginResponse](apiClient, types.UpdateSelfEndpoint, request, client.AddAuthHeader)
	if err != nil {
		console.Error(err.Error())
		return
	}
	store.SetUser(result.User)
}
