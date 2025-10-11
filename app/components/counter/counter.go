//go:build js && wasm

package counter

import (
	"fmt"

	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/breakpoint"
	"github.com/gofred-io/gofred/foundation/button"
	"github.com/gofred-io/gofred/foundation/column"
	"github.com/gofred-io/gofred/foundation/container"
	"github.com/gofred-io/gofred/foundation/row"
	"github.com/gofred-io/gofred/foundation/spacer"
	"github.com/gofred-io/gofred/foundation/text"
	"github.com/gofred-io/gofred/hooks"
	"github.com/gofred-io/gofred/listenable"
	"github.com/gofred-io/gofred/options/spacing"
	"github.com/misleb/mego2/app/client"
)

var (
	count, setCount    = hooks.UseState(0)
	errorMsg, setError = hooks.UseState("")
)

func Get() application.BaseWidget {
	return container.New(
		column.New(
			[]application.BaseWidget{
				listenable.Builder(errorMsg, func() application.BaseWidget {
					return text.New(
						errorMsg.Value(),
						text.FontSize(24),
						text.FontColor("#FF0000"),
					)
				}),
				text.New(
					"Counter App",
					text.FontSize(24),
					text.FontColor("#1F2937"),
					text.FontWeight("700"),
				),
				listenable.Builder(count, func() application.BaseWidget {
					return text.New(
						fmt.Sprintf("Count: %d", count.Value()),
						text.FontSize(18),
						text.FontColor("#2B799B"),
						text.FontWeight("600"),
					)
				}),
				spacer.New(spacer.Height(16)),
				row.New(
					[]application.BaseWidget{
						button.New(
							text.New("Decrease", text.FontColor("#FFFFFF")),
							button.BackgroundColor("#EF4444"),
							button.OnClick(decreaseCount),
						),
						spacer.New(spacer.Width(16)),
						button.New(
							text.New("Increase", text.FontColor("#FFFFFF")),
							button.BackgroundColor("#10B981"),
							button.OnClick(increaseCount),
						),
					},
					row.Gap(16),
				),
			},
			column.Gap(16),
		),
		container.Padding(breakpoint.All(spacing.All(32))),
		container.BackgroundColor("#FFFFFF"),
	)
}

func increaseCount(this application.BaseWidget, e application.Event) {
	apiClient := client.GetInstance()
	result, err := apiClient.Increment(count.Value())
	if err != nil {
		setError(fmt.Sprintf("Error: %v", err))
		return
	}
	setError("")
	setCount(result.Result)
}

func decreaseCount(this application.BaseWidget, e application.Event) {
	setCount(count.Value() - 1)
}
