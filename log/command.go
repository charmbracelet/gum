package log

import (
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

// Run is the command-line interface for logging text.
func (o Options) Run() error {
	l := log.New(os.Stderr)

	if o.File != "" {
		f, err := os.OpenFile(o.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		defer f.Close() // nolint: errcheck
		l.SetOutput(f)
	}

	l.SetPrefix(o.Prefix)
	l.SetLevel(-math.MaxInt32) // log all levels
	l.SetReportTimestamp(o.Time)
	l.SetTimeFormat(o.TimeFormat)

	st := log.DefaultStyles()
	st.Levels[log.DebugLevel] = o.LevelDebugStyle.ToLipgloss().
		Inline(true).
		SetString(strings.ToUpper(log.DebugLevel.String()))
	st.Levels[log.InfoLevel] = o.LevelInfoStyle.ToLipgloss().
		Inline(true).
		SetString(strings.ToUpper(log.InfoLevel.String()))
	st.Levels[log.WarnLevel] = o.LevelWarnStyle.ToLipgloss().
		Inline(true).
		SetString(strings.ToUpper(log.WarnLevel.String()))
	st.Levels[log.ErrorLevel] = o.LevelErrorStyle.ToLipgloss().
		Inline(true).
		SetString(strings.ToUpper(log.ErrorLevel.String()))
	st.Levels[log.FatalLevel] = o.LevelFatalStyle.ToLipgloss().
		Inline(true).
		SetString(strings.ToUpper(log.FatalLevel.String()))
	st.Timestamp = o.TimeStyle.ToLipgloss().
		Inline(true)
	st.Prefix = o.PrefixStyle.ToLipgloss().
		Inline(true)
	st.Message = o.MessageStyle.ToLipgloss().
		Inline(true)
	st.Key = o.KeyStyle.ToLipgloss().
		Inline(true)
	st.Value = o.ValueStyle.ToLipgloss().
		Inline(true)
	st.Separator = o.SeparatorStyle.ToLipgloss().
		Inline(true)

	l.SetStyles(st)

	switch o.Formatter {
	case "json":
		l.SetFormatter(log.JSONFormatter)
	case "logfmt":
		l.SetFormatter(log.LogfmtFormatter)
	case "text":
		l.SetFormatter(log.TextFormatter)
	}

	var arg0 string
	var args []interface{}
	if len(o.Text) > 0 {
		arg0 = o.Text[0]
	}

	if len(o.Text) > 1 {
		args = make([]interface{}, len(o.Text[1:]))
		for i, arg := range o.Text[1:] {
			args[i] = arg
		}
	}

	logger := map[string]logger{
		"none":  {printf: l.Printf, print: l.Print},
		"debug": {printf: l.Debugf, print: l.Debug},
		"info":  {printf: l.Infof, print: l.Info},
		"warn":  {printf: l.Warnf, print: l.Warn},
		"error": {printf: l.Errorf, print: l.Error},
		"fatal": {printf: l.Fatalf, print: l.Fatal},
	}[o.Level]

	if o.Format {
		logger.printf(arg0, args...)
	} else if o.Structured {
		logger.print(arg0, args...)
	} else {
		logger.print(strings.Join(o.Text, " "))
	}

	return nil
}

type logger struct {
	printf func(string, ...interface{})
	print  func(interface{}, ...interface{})
}
