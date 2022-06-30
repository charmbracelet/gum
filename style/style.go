package style

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/coral"
)

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
				str, err = readStdin()
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

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return "", nil
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}
