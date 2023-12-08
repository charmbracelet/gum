package progress

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

type barFormatter struct {
	pbar    progress.Model
	numBars int
	tplstr  string
}

var barPlaceholderRe = regexp.MustCompile(`{\s*Bar\s*}`)
var nonBarPlaceholderRe = regexp.MustCompile(`{\s*(Title|Elapsed|Iter|Avg|Pct|Eta|Remaining|Limit)\s*}`)

func newBarFormatter(tplstr string, barColor string) *barFormatter {
	var bar progress.Model
	if barColor != "" {
		bar = progress.New(progress.WithoutPercentage(), progress.WithSolidFill(barColor))
	} else {
		bar = progress.New(progress.WithoutPercentage())
	}
	barfmt := &barFormatter{
		pbar:    bar,
		tplstr:  tplstr,
		numBars: len(barPlaceholderRe.FindAllString(tplstr, -1)),
	}
	return barfmt
}

func (self *barFormatter) Render(info *barInfo, maxWidth int) string {
	rendered := nonBarPlaceholderRe.ReplaceAllStringFunc(self.tplstr, func(s string) string {
		switch strings.TrimSpace(s[1 : len(s)-1]) {
		case "Title":
			return info.title
		case "Iter":
			return fmt.Sprint(info.iter)
		case "Limit":
			if info.limit == 0 {
				return s
			}
			return fmt.Sprint(info.limit)
		case "Elapsed":
			return info.Elapsed().String()
		case "Pct":
			if info.limit == 0 {
				return s
			}
			return fmt.Sprintf("%d%%", info.Pct())
		case "Avg":
			return info.Avg().Round(time.Second).String()
		case "Remaining":
			if info.limit == 0 {
				return s
			}
			return info.Eta().Round(time.Second).String()
		case "Eta":
			if info.limit == 0 {
				return s
			}
			return time.Now().Add(info.Eta()).Format(time.TimeOnly)
		default:
			return ""
		}
	})

	if info.limit > 0 && self.numBars > 0 {
		self.pbar.Width = max(0, (maxWidth-lipgloss.Width(rendered))/int(self.numBars))
		bar := self.pbar.ViewAs(safeDivide(float64(info.iter), float64(info.limit)))
		rendered = barPlaceholderRe.ReplaceAllLiteralString(rendered, bar)
	}
	return rendered
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}
