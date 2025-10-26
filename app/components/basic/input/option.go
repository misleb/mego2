//go:build js && wasm

package input

import (
	"github.com/gofred-io/gofred/application"
	"github.com/gofred-io/gofred/options"
)

// InputType represents the possible HTML input types
type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypePassword InputType = "password"
	InputTypeEmail    InputType = "email"
	InputTypeNumber   InputType = "number"
	InputTypeTel      InputType = "tel"
	InputTypeURL      InputType = "url"
	InputTypeSearch   InputType = "search"
	InputTypeDate     InputType = "date"
	InputTypeDateTime InputType = "datetime-local"
	InputTypeTime     InputType = "time"
	InputTypeMonth    InputType = "month"
	InputTypeWeek     InputType = "week"
	InputTypeColor    InputType = "color"
	InputTypeCheckbox InputType = "checkbox"
	InputTypeRadio    InputType = "radio"
	InputTypeRange    InputType = "range"
	InputTypeFile     InputType = "file"
	InputTypeHidden   InputType = "hidden"
	InputTypeImage    InputType = "image"
	InputTypeSubmit   InputType = "submit"
	InputTypeButton   InputType = "button"
	InputTypeReset    InputType = "reset"
)

type Option options.OptionWrapper

func Placeholder(ph string) Option {
	return func() options.Option {
		return func(widget application.BaseWidget) {
			widget.SetAttribute("placeholder", ph)
		}
	}
}

func Type(t InputType) Option {
	return func() options.Option {
		return func(widget application.BaseWidget) {
			widget.SetAttribute("type", string(t))
		}
	}
}
