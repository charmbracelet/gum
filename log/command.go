package log

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Run is the command-line interface for logging text.
func (o Options) Run() error {
	l := log.New(os.Stderr)

	if o.File != "" {
		f, err := os.OpenFile(o.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}

		defer f.Close() //nolint:errcheck
		l.SetOutput(f)
	}

	l.SetPrefix(o.Prefix)
	l.SetLevel(-math.MaxInt32) // log all levels
	l.SetReportTimestamp(o.Time)

	timeFormats := map[string]string{
		"layout":      "01/02 03:04:05PM '06 -0700",
		"ansic":       "Mon Jan _2 15:04:05 2006",
		"unixdate":    "Mon Jan _2 15:04:05 MST 2006",
		"rubydate":    "Mon Jan 02 15:04:05 -0700 2006",
		"rfc822":      "02 Jan 06 15:04 MST",
		"rfc822z":     "02 Jan 06 15:04 -0700",
		"rfc850":      "Monday, 02-Jan-06 15:04:05 MST",
		"rfc1123":     "Mon, 02 Jan 2006 15:04:05 MST",
		"rfc1123z":    "Mon, 02 Jan 2006 15:04:05 -0700",
		"rfc3339":     "2006-01-02T15:04:05Z07:00",
		"rfc3339nano": "2006-01-02T15:04:05.999999999Z07:00",
		"kitchen":     "3:04PM",
		"stamp":       "Jan _2 15:04:05",
		"stampmilli":  "Jan _2 15:04:05.000",
		"stampmicro":  "Jan _2 15:04:05.000000",
		"stampnano":   "Jan _2 15:04:05.000000000",
		"datetime":    "2006-01-02 15:04:05",
		"dateonly":    "2006-01-02",
		"timeonly":    "15:04:05",
	}

	tf, ok := timeFormats[strings.ToLower(o.TimeFormat)]
	if ok {
		l.SetTimeFormat(tf)
	} else {
		l.SetTimeFormat(o.TimeFormat)
	}

	st := log.DefaultStyles()
	defaultColors := map[log.Level]lipgloss.Color{
		log.DebugLevel: lipgloss.Color("63"),
		log.InfoLevel:  lipgloss.Color("83"),
		log.WarnLevel:  lipgloss.Color("192"),
		log.ErrorLevel: lipgloss.Color("204"),
		log.FatalLevel: lipgloss.Color("134"),
	}

	lvlStyle := o.LevelStyle.ToLipgloss()
	if lvlStyle.GetForeground() == lipgloss.Color("") {
		lvlStyle = lvlStyle.Foreground(defaultColors[levelToLog[o.Level]])
	}

	lvl := levelToLog[o.Level]
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
