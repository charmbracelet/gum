package style

// Options is the customization options for the style command.
type Options struct {
	Background       string `help:"Background color"`
	Foreground       string `help:"Foreground color"`
	BorderBackground string `help:"Border background color"`
	BorderForeground string `help:"Border foreground color"`
	Align            string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left"`
	Border           string `help:"Border style to apply" enum:"none,hidden,normal,rounded,thick,double" default:"none"`
	Height           int    `help:"Height of output"`
	Width            int    `help:"Width of output"`
	Margin           string `help:"Margin to apply around the text."`
	Padding          string `help:"Padding to apply around the text."`
	Bold             bool   `help:"Apply bold formatting"`
	Faint            bool   `help:"Apply faint formatting"`
	Italic           bool   `help:"Apply italic formatting"`
	Strikethrough    bool   `help:"Apply strikethrough formatting"`
}
