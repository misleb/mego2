//go:build js && wasm

package home

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/foundation/column"
	"github.com/gofred-io/gofred/foundation/router"
	"github.com/gofred-io/gofred/listenable"
	"github.com/misleb/mego2/app/components/account"
	"github.com/misleb/mego2/app/components/counter"
	"github.com/misleb/mego2/app/components/header"
	"github.com/misleb/mego2/app/components/login"
	"github.com/misleb/mego2/app/store"
)

func New(params router.RouteParams) application.BaseWidget {
	return column.New(
		[]application.BaseWidget{
			header.Get(),
			listenable.Builder(store.AppStoreListenable(), func() application.BaseWidget {
				user := store.GetUser()
				if user == nil {
					return login.Get()
				} else if user.SetPassword {
					return account.Get()
				} else {
					return counter.Get()
				}
			}),
		},
	)
}
