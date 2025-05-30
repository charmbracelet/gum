// Package log the log command.
package log

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Run is the command-line interface for logging text.
func (o Options) Run() error {
	l := log.New(os.Stderr)

	if o.File != "" {
		f, err := os.OpenFile(o.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm) //nolint:gosec
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}

		defer f.Close() //nolint:errcheck
		l.SetOutput(f)
	}

	l.SetPrefix(o.Prefix)
	l.SetLevel(-math.MaxInt32) // log all levels
	l.SetReportTimestamp(o.Time != "")
	if o.MinLevel != "" {
		lvl, err := log.ParseLevel(o.MinLevel)
		if err != nil {
			return err //nolint:wrapcheck
		}
		l.SetLevel(lvl)
	}

	timeFormats := map[string]string{
		"layout":      time.Layout,
		"ansic":       time.ANSIC,
		"unixdate":    time.UnixDate,
		"rubydate":    time.RubyDate,
		"rfc822":      time.RFC822,
		"rfc822z":     time.RFC822Z,
		"rfc850":      time.RFC850,
		"rfc1123":     time.RFC1123,
		"rfc1123z":    time.RFC1123Z,
		"rfc3339":     time.RFC3339,
		"rfc3339nano": time.RFC3339Nano,
		"kitchen":     time.Kitchen,
		"stamp":       time.Stamp,
		"stampmilli":  time.StampMilli,
		"stampmicro":  time.StampMicro,
		"stampnano":   time.StampNano,
		"datetime":    time.DateTime,
		"dateonly":    time.DateOnly,
		"timeonly":    time.TimeOnly,
	}

	tf, ok := timeFormats[strings.ToLower(o.Time)]
	if ok {
		l.SetTimeFormat(tf)
	} else {
		l.SetTimeFormat(o.Time)
	}

	st := log.DefaultStyles()
	lvl := levelToLog[o.Level]
	lvlStyle := o.LevelStyle.ToLipgloss()
	if lvlStyle.GetForeground() == lipgloss.Color("") {
		lvlStyle = lvlStyle.Foreground(st.Levels[lvl].GetForeground())
	}

	st.Levels[lvl] = lvlStyle.
		SetString(strings.ToUpper(lvl.String())).
		Inline(true)

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

var levelToLog = map[string]log.Level{
	"none":  log.Level(math.MaxInt32),
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
}
