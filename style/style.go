package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/coral"
)

type options struct {
	// Colors
	background       *string
	foreground       *string
	borderForeground *string
	borderBackground *string

	// Decoration
	border *string

	// Layout
	align   *string
	height  *int
	width   *int
	margin  *string
	padding *string

	// Styles
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
		Args:  coral.ExactArgs(1),
		Run: func(cmd *coral.Command, args []string) {
			s := lipgloss.NewStyle()

			// Color
			s = s.Background(lipgloss.Color(*opts.background))
			s = s.Foreground(lipgloss.Color(*opts.foreground))
			s = s.BorderForeground(lipgloss.Color(*opts.borderForeground))
			s = s.BorderBackground(lipgloss.Color(*opts.borderBackground))

			// Layout
			s = s.Align(alignment[*opts.align])
			s = s.Height(*opts.height)
			s = s.Width(*opts.width)
			s = s.Padding(parsePadding(*opts.padding))
			s = s.Margin(parsePadding(*opts.margin))

			// Decoration
			s = s.Border(border[*opts.border])

			// Style
			s = s.Bold(*opts.bold)
			s = s.Faint(*opts.faint)
			s = s.Italic(*opts.italic)
			s = s.Strikethrough(*opts.strikethrough)

			fmt.Println(s.Render(args[0]))
		},
	}

	f := cmd.Flags()
	opts = options{
		// Colors
		background:       f.String("background", "", "Background"),
		foreground:       f.String("foreground", "", "Foreground"),
		borderBackground: f.String("border-background", "", "Border Background"),
		borderForeground: f.String("border-foreground", "", "Border Foreground"),

		// Decoration
		border: f.String("border", "none", "Border"),

		// Layout
		align:   f.String("align", "left", "Alignment"),
		height:  f.Int("height", 0, "Height"),
		width:   f.Int("width", 0, "Width"),
		padding: f.String("padding", "0 0 0 0", "Padding"),
		margin:  f.String("margin", "0 0 0 0", "Margin"),

		// Styles
		bold:          f.Bool("bold", false, "Bold"),
		faint:         f.Bool("faint", false, "Faint"),
		italic:        f.Bool("italic", false, "Italic"),
		strikethrough: f.Bool("strikethrough", false, "Strikethrough"),
	}

	return cmd
}
