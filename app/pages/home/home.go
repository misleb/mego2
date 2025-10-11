//go:build js && wasm

package home

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/foundation/column"
	"github.com/gofred-io/gofred/foundation/router"
	"github.com/misleb/mego2/app/components/counter"
	"github.com/misleb/mego2/app/components/header"
)

func New(params router.RouteParams) application.BaseWidget {
	return column.New(
		[]application.BaseWidget{
			header.Get(),
			counter.Get(),
		},
	)
}
