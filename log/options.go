package log

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/style"
)

// Options is the set of options that can configure a join.
type Options struct {
	Text []string `arg:"" help:"Text to log"`

	File       string `short:"o" help:"Log to file"`
	Format     bool   `short:"f" help:"Format message using printf" xor:"format,structured"`
	Formatter  string `help:"The log formatter to use" enum:"json,logfmt,text" default:"text"`
	Level      string `help:"The log level to use" enum:"none,debug,info,warn,error,fatal" default:"none"`
	Prefix     string `help:"Prefix to print before the message"`
	Structured bool   `short:"s" help:"Use structured logging" xor:"format,structured"`
	Time       bool   `help:"Whether to print the time"`
	TimeFormat string `help:"The time format to use" default:"2006/01/02 15:04:05"`

	LevelDebugStyle style.Styles `embed:"" prefix:"level.debug." help:"The style of the debug level" set:"defaultBold=true" set:"defaultForeground=63" envprefix:"GUM_LOG_LEVEL_DEBUG_"`  //nolint:staticcheck
	LevelInfoStyle  style.Styles `embed:"" prefix:"level.info." help:"The style of the info level" set:"defaultBold=true" set:"defaultForeground=83" envprefix:"GUM_LOG_LEVEL_INFO_"`     //nolint:staticcheck
	LevelWarnStyle  style.Styles `embed:"" prefix:"level.warn." help:"The style of the warn level" set:"defaultBold=true" set:"defaultForeground=192" envprefix:"GUM_LOG_LEVEL_WARN_"`    //nolint:staticcheck
	LevelErrorStyle style.Styles `embed:"" prefix:"level.error." help:"The style of the error level" set:"defaultBold=true" set:"defaultForeground=204" envprefix:"GUM_LOG_LEVEL_ERROR_"` //nolint:staticcheck
	LevelFatalStyle style.Styles `embed:"" prefix:"level.fatal." help:"The style of the fatal level" set:"defaultBold=true" set:"defaultForeground=134" envprefix:"GUM_LOG_LEVEL_FATAL_"` //nolint:staticcheck
	TimeStyle       style.Styles `embed:"" prefix:"time." help:"The style of the time" envprefix:"GUM_LOG_TIME_"`
	PrefixStyle     style.Styles `embed:"" prefix:"prefix." help:"The style of the prefix" set:"defaultBold=true" set:"defaultFaint=true" envprefix:"GUM_LOG_PREFIX_"` //nolint:staticcheck
	MessageStyle    style.Styles `embed:"" prefix:"message." help:"The style of the message" envprefix:"GUM_LOG_MESSAGE_"`
	KeyStyle        style.Styles `embed:"" prefix:"key." help:"The style of the key" set:"defaultFaint=true" envprefix:"GUM_LOG_KEY_"`
	ValueStyle      style.Styles `embed:"" prefix:"value." help:"The style of the value" envprefix:"GUM_LOG_VALUE_"`
	SeparatorStyle  style.Styles `embed:"" prefix:"separator." help:"The style of the separator" set:"defaultFaint=true" envprefix:"GUM_LOG_SEPARATOR_"`
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}
