//go:build js && wasm

package store

import (
	"github.com/gofred-io/gofred/hooks"
	"github.com/gofred-io/gofred/listenable"
	"github.com/misleb/mego2/shared"
)

type Notifcation struct {
	Message string
	Type    string
}

type AppStore struct {
	User          *shared.User
	Theme         string
	Notifications []Notifcation
}

var (
	appStore, setAppStore = hooks.UseState(AppStore{
		Theme:         "dark",
		Notifications: []Notifcation{},
	})
)

func SetUser(user *shared.User) {
	store := appStore.Value()
	store.User = user
	setAppStore(store)
}

func GetUser() *shared.User {
	return appStore.Value().User
}

func AppStoreListenable() listenable.Listenable[AppStore] {
	return appStore
}
