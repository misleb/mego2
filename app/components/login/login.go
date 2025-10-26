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
	"github.com/gofred-io/gofred/hooks"
	"github.com/gofred-io/gofred/options/spacing"
	"github.com/misleb/mego2/app/client"
	"github.com/misleb/mego2/app/components/basic/input"
	"github.com/misleb/mego2/app/store"
	"github.com/misleb/mego2/shared"
)

var (
	IsLoggedIn, SetLoggedIn = hooks.UseState(false)
	userInput               *input.Input
	passInput               *input.Input
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
			},
		),
		container.Padding(breakpoint.All(spacing.All(32))),
	)
}

func doLogin(this application.BaseWidget, e application.Event) {
	client := client.GetInstance()
	result, err := client.Login(userInput.GetValue(), passInput.GetValue())
	if err != nil {
		console.Error(result.Error)
		return
	}
	store.SetUser(&shared.User{Name: userInput.GetValue(), Token: result.Token})
}
