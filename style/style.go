package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/seashell/internal/stdin"
	"github.com/muesli/coral"
)

// alignment maps strings to lipgloss.Position's
var alignment = map[string]lipgloss.Position{
	"center": lipgloss.Center,
	"left":   lipgloss.Left,
	"top":    lipgloss.Top,
	"bottom": lipgloss.Bottom,
	"right":  lipgloss.Right,
}

// border maps strings to lipgloss.Border's
var border map[string]lipgloss.Border = map[string]lipgloss.Border{
	"double":  lipgloss.DoubleBorder(),
	"hidden":  lipgloss.HiddenBorder(),
	"none":    {},
	"normal":  lipgloss.NormalBorder(),
	"rounded": lipgloss.RoundedBorder(),
	"thick":   lipgloss.ThickBorder(),
}

type options struct {
	// Colors
	background       *string
	foreground       *string
	borderBackground *string
	borderForeground *string
	// Layout
	align   *string
	border  *string
	margin  *string
	padding *string
	height  *int
	width   *int
	// Format
	bold          *bool
	faint         *bool
	italic        *bool
	strikethrough *bool
}

// Cmd returns the style command
func Cmd() *coral.Command {
	var opts options

	var cmd = &coral.Command{
		Use:   "style",
		Short: "Style styles a given string",
		Args:  coral.ArbitraryArgs,
		Run: func(cmd *coral.Command, args []string) {
			var str string
			var err error

			if len(args) <= 0 {
				// No arguments are passed, let's check stdin
				str, err = stdin.Read()
				if err != nil || str == "" {
					// No stdin, let's display the help
					cmd.Help()
					return
				}
			} else {
				str = strings.Join(args, " ")
			}

			fmt.Println(lipgloss.NewStyle().
				// Colors
				Foreground(lipgloss.Color(*opts.foreground)).
				Background(lipgloss.Color(*opts.background)).
				BorderBackground(lipgloss.Color(*opts.borderBackground)).
				BorderForeground(lipgloss.Color(*opts.borderForeground)).
				// Layout
				Align(alignment[*opts.align]).
				Bold(*opts.bold).
				Border(border[*opts.border]).
				Margin(parseMargin(*opts.margin)).
				Padding(parsePadding(*opts.padding)).
				Height(*opts.height).
				Width(*opts.width).
				// Format
				Faint(*opts.faint).
				Italic(*opts.italic).
				Strikethrough(*opts.strikethrough).
				// Render the string,
				// with the styling specified through the flag options.
				Render(str))

		},
	}

	f := cmd.Flags()
	opts = options{
		// Colors
		background:       f.StringP("background", "b", "", "The background color to apply"),
		foreground:       f.StringP("foreground", "f", "", "The foreground color to apply"),
		borderBackground: f.String("border-background", "", "The border background color to apply"),
		borderForeground: f.String("border-foreground", "", "The border foreground color to apply"),

		// Layout
		align:   f.StringP("align", "a", "left", "The text alignment (left, center, right, bottom, middle, top)"),
		border:  f.String("border", "none", "The border style to apply (none, hidden, normal, rounded, thick, double)"),
		height:  f.IntP("height", "h", 0, "The height the output should take up"),
		width:   f.IntP("width", "w", 0, "The width the output should take up"),
		margin:  f.StringP("margin", "m", "0 0 0 0", "Margin to apply around the text."),
		padding: f.StringP("padding", "p", "0 0 0 0", "Padding to apply around the text."),

		// Format
		bold:          f.Bool("bold", false, "Whether to apply bold formatting"),
		faint:         f.Bool("faint", false, "Whether to apply faint formatting"),
		italic:        f.BoolP("italic", "i", false, "Whether to apply italic formatting"),
		strikethrough: f.BoolP("strikethrough", "s", false, "Whether to apply strikethrough formatting"),
	}

	return cmd
}

// parsePadding parses 1 - 4 integers from a string and returns them in a top,
// right, bottom, left order for use in the lipgloss.Padding() method
func parsePadding(s string) (int, int, int, int) {
	var ints []int

	tokens := strings.Split(s, " ")

	// All tokens must be an integer
	for _, token := range tokens {
		parsed, err := strconv.Atoi(token)
		if err != nil {
			return 0, 0, 0, 0
		}
		ints = append(ints, parsed)
	}

	if len(tokens) == 1 {
		return ints[0], ints[0], ints[0], ints[0]
	}

	if len(tokens) == 2 {
		return ints[0], ints[1], ints[0], ints[1]
	}

	if len(tokens) == 4 {
		return ints[0], ints[1], ints[2], ints[3]
	}

	return 0, 0, 0, 0
}

// parseMargin is an alias for parsePadding since they involve the same logic
// to parse integers to the same format.
var parseMargin = parsePadding
