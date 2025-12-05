//go:build js && wasm

package store

import (
	"github.com/gofred-io/gofred/hooks"
	"github.com/gofred-io/gofred/listenable"
	"github.com/misleb/mego2/shared/types"
)

type Notifcation struct {
	Message string
	Type    string
}

type AppStore struct {
	User          *types.User
	Theme         string
	Notifications []Notifcation
	AuthToken     string
}

var (
	appStore, setAppStore = hooks.UseState(AppStore{
		Theme:         "dark",
		Notifications: []Notifcation{},
	})
)

func SetUser(user *types.User) {
	store := appStore.Value()
	store.User = user
	if user.CurrentToken != "" {
		store.AuthToken = user.CurrentToken
	}
	setAppStore(store)
}

func GetUser() *types.User {
	return appStore.Value().User
}

func AppStoreListenable() listenable.Listenable[AppStore] {
	return appStore
}

func GetAuthToken() string {
	return appStore.Value().AuthToken
}
