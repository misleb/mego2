//go:build js && wasm

package main

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/foundation/router"
	"github.com/misleb/mego2/app/pages/home"
)

func main() {
	application.Run(router.New(
		router.Route("/", home.New),
	))
}
