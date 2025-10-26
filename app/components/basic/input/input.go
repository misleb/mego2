//go:build js && wasm

package input

import "github.com/gofred-io/gofred/application"

type Input struct {
	application.BaseWidget
}

func New(defaultValue string, options ...Option) *Input {
	input := &Input{
		BaseWidget: application.New("input"),
	}

	input.SetValue(defaultValue)

	for _, option := range options {
		option()(input.BaseWidget)
	}

	return input
}

func (i *Input) GetValue() string {
	return i.Get("value").String()
}

func (i *Input) SetValue(s string) {
	i.Set("value", s)
}
