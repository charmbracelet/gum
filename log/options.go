package log

import (
	"github.com/charmbracelet/gum/style"
)

// Options is the set of options that can configure a join.
type Options struct {
	Text []string `arg:"" help:"Text to log"`

	File       string `short:"o" help:"Log to file"`
	Format     bool   `short:"f" help:"Format message using printf" xor:"format,structured"`
	Formatter  string `help:"The log formatter to use" enum:"json,logfmt,text" default:"text"`
	Level      string `short:"l" help:"The log level to use" enum:"none,debug,info,warn,error,fatal" default:"none"`
	Prefix     string `help:"Prefix to print before the message"`
	Structured bool   `short:"s" help:"Use structured logging" xor:"format,structured"`
	Time       string `short:"t" help:"The time format to use (kitchen, layout, ansic, rfc822, etc...)" default:""`

	MinLevel string `help:"Minimal level to show" default:"" env:"GUM_LOG_LEVEL"`

	LevelStyle     style.Styles `embed:"" prefix:"level." help:"The style of the level being used" set:"defaultBold=true" envprefix:"GUM_LOG_LEVEL_"`
	TimeStyle      style.Styles `embed:"" prefix:"time." help:"The style of the time" envprefix:"GUM_LOG_TIME_"`
	PrefixStyle    style.Styles `embed:"" prefix:"prefix." help:"The style of the prefix" set:"defaultBold=true" set:"defaultFaint=true" envprefix:"GUM_LOG_PREFIX_"` //nolint:staticcheck
	MessageStyle   style.Styles `embed:"" prefix:"message." help:"The style of the message" envprefix:"GUM_LOG_MESSAGE_"`
	KeyStyle       style.Styles `embed:"" prefix:"key." help:"The style of the key" set:"defaultFaint=true" envprefix:"GUM_LOG_KEY_"`
	ValueStyle     style.Styles `embed:"" prefix:"value." help:"The style of the value" envprefix:"GUM_LOG_VALUE_"`
	SeparatorStyle style.Styles `embed:"" prefix:"separator." help:"The style of the separator" set:"defaultFaint=true" envprefix:"GUM_LOG_SEPARATOR_"`
}
