package timeout

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Options Common Timeout Option.
type Options struct {
	Timeout time.Duration `help:"Timeout until command exits" default:"0" env:"GUM_CONFIRM_TIMEOUT"`
}

// HasTimeout checks for a given timeout parameter.
func (o Options) HasTimeout() (hasTimeout bool) {
	return o.Timeout > 0
}

// Tick interval.
const tickInterval = time.Second

// TickTimeoutMsg will be dispatched for every tick.
// Containing current timeout value
// and optional parameter to be used when handling the timeout msg.
type TickTimeoutMsg struct {
	TimeoutValue time.Duration
	Data         interface{}
}

// Init Start Timeout ticker using with timeout in seconds and optional data.
func Init(timeout time.Duration, data interface{}) tea.Cmd {
	if timeout > 0 {
		return Tick(timeout, data)
	}
	return nil
}

// Start ticker.
func Tick(timeoutValue time.Duration, data interface{}) tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		// every tick checks if the timeout needs to be decremented
		// and send as message
		if timeoutValue >= 0 {
			timeoutValue -= tickInterval
			return TickTimeoutMsg{
				TimeoutValue: timeoutValue,
				Data:         data,
			}
		}
		return nil
	})
}

// Str produce Timeout String to be rendered.
func Str(timeout time.Duration) string {
	return fmt.Sprintf(" (%d)", max(0, int(timeout.Seconds())))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
