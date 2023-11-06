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

	switch o.Level {
	case "none":
		if o.Format {
			l.Printf(arg0, args...)
		} else if o.Structured {
			l.Print(arg0, args...)
		} else {
			l.Print(strings.Join(o.Text, " "))
		}
	case "debug":
		if o.Format {
			l.Debugf(arg0, args...)
		} else if o.Structured {
			l.Debug(arg0, args...)
		} else {
			l.Debug(strings.Join(o.Text, " "))
		}
	case "info":
		if o.Format {
			l.Infof(arg0, args...)
		} else if o.Structured {
			l.Info(arg0, args...)
		} else {
			l.Info(strings.Join(o.Text, " "))
		}
	case "warn":
		if o.Format {
			l.Warnf(arg0, args...)
		} else if o.Structured {
			l.Warn(arg0, args...)
		} else {
			l.Warn(strings.Join(o.Text, " "))
		}
	case "error":
		if o.Format {
			l.Errorf(arg0, args...)
		} else if o.Structured {
			l.Error(arg0, args...)
		} else {
			l.Error(strings.Join(o.Text, " "))
		}
	case "fatal":
		if o.Format {
			l.Fatalf(arg0, args...)
		} else if o.Structured {
			l.Fatal(arg0, args...)
		} else {
			l.Fatal(strings.Join(o.Text, " "))
		}
	}

	return nil
}
