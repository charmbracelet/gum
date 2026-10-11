package completion

import (
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

func TestFish(t *testing.T) {
	var cli struct {
		Choose struct {
			Options []string `arg:"" optional:""`
			Cursor  string   `help:"Cursor"`
		} `cmd:"" help:"Choose an option"`
		Completion Completion `cmd:"" hidden:""`
	}
	var out strings.Builder
	parser := kong.Must(&cli, kong.Name("gum"), kong.Writers(&out, &out))
	ctx, err := parser.Parse([]string{"completion", "fish"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ctx.Run(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(out.String(), "\n")

	for name, tt := range map[string]struct {
		line string
		want bool
	}{
		"files after a subcommand": {
			line: "complete -c gum -f",
			want: false,
		},
		"no files in place of a subcommand": {
			line: "complete -c gum -f -n '__fish_use_subcommand' -a choose -d 'Choose an option'",
			want: true,
		},
		"no files for a flag value": {
			line: `complete -c gum -f -n '__fish_seen_subcommand_from choose' -x -l cursor -d "Cursor"`,
			want: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := slices.Contains(lines, tt.line); got != tt.want {
				t.Errorf("expected %v, got %v for %q", tt.want, got, tt.line)
			}
		})
	}
}
