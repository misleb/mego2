//go:build js && wasm

package header

import (
	"fmt"

	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/breakpoint"
	"github.com/gofred-io/gofred/foundation/container"
	"github.com/gofred-io/gofred/foundation/row"
	"github.com/gofred-io/gofred/foundation/text"
	"github.com/gofred-io/gofred/listenable"
	"github.com/gofred-io/gofred/options/spacing"
	"github.com/gofred-io/gofred/theme"
	"github.com/misleb/mego2/app/store"
)

func Get() application.BaseWidget {
	return container.New(
		row.New(
			[]application.BaseWidget{
				listenable.Builder(store.AppStoreListenable(), func() application.BaseWidget {
					user := store.GetUser()
					if user != nil {
						return text.New(fmt.Sprintf("Welcome %s", user.Name))
					}
					return text.New("Header", text.FontSize(24), text.FontColor("#000000"))
				}),
			},
			row.Gap(0),
			row.Flex(1),
			row.CrossAxisAlignment(theme.AxisAlignmentTypeCenter),
		),
		container.Height(breakpoint.All(72)),
		container.BackgroundColor("#FFFFFF"),
		container.Padding(
			breakpoint.All(spacing.Axis(32, 0)),
			breakpoint.XS(spacing.Axis(0, 0)),
			breakpoint.SM(spacing.Axis(8, 0)),
			breakpoint.MD(spacing.Axis(16, 0)),
			breakpoint.LG(spacing.Axis(24, 0)),
		),
		container.BorderColor("#E5E7EB"),
		container.BorderWidth(spacing.Spacing{0, 0, 1, 0}),
		container.BorderStyle(theme.BorderStyleTypeSolid),
		// Add subtle shadow for depth
		container.BoxShadow("0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)"),
	)
}
