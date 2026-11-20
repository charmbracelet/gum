package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestExtendedFilterCLI(t *testing.T) {
	const helper = "GUM_TEST_EXTENDED_FILTER"
	if os.Getenv(helper) == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				os.Args = append([]string{"gum"}, os.Args[i+1:]...)
				main()
				os.Exit(0)
			}
		}
		t.Fatal("missing CLI arguments")
	}
	for _, tc := range []struct {
		name  string
		args  []string
		env   []string
		stdin string
		want  string
	}{
		{
			name: "flag and initial value", args: []string{"--extended", "--value", "src 'model"},
			stdin: "docs/model.md\nsrc/model.go\nsrc/main.go", want: "src/model.go\n",
		},
		{
			name: "environment variable", args: []string{"--value", "src 'model"},
			env:   []string{"GUM_FILTER_EXTENDED=true"},
			stdin: "docs/model.md\nsrc/model.go\nsrc/main.go", want: "src/model.go\n",
		},
		{
			name: "exact filtering", args: []string{"--extended", "--no-fuzzy", "--value", "src model"},
			stdin: "src/model.go\nsrc/m_o_d_e_l.go", want: "src/model.go\n",
		},
		{
			name: "escaped phrase", args: []string{"--extended", "--value", `'user\ model`},
			stdin: "user model.go\nuser_model.go", want: "user model.go\n",
		},
		{
			name: "Unicode case folding", args: []string{"--extended", "--value", "'kelvin 'go"},
			stdin: "Kelvin.go\nKelvin.md", want: "Kelvin.go\n",
		},
		{
			name: "styled options", args: []string{"--extended", "--value", "src 'model"},
			stdin: "\x1b[31msrc/model.go\x1b[0m\ndocs/model.md", want: "src/model.go\n",
		},
		{
			name: "argument options", args: []string{"--extended", "--value", "src 'model", "src/model.go", "docs/model.md"},
			want: "src/model.go\n",
		},
		{
			name: "legacy apostrophe remains literal", args: []string{"--value", "'model"},
			stdin: "'model.go\nmodel.go", want: "'model.go\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			args := []string{"-test.run=^TestExtendedFilterCLI$", "--", "filter", "--select-if-one"}
			cmd := exec.CommandContext(ctx, os.Args[0], append(args, tc.args...)...)
			// Keep user customization out of the CLI compatibility assertions.
			for _, env := range os.Environ() {
				if !strings.HasPrefix(env, "GUM_FILTER_") && !strings.HasPrefix(env, helper+"=") {
					cmd.Env = append(cmd.Env, env)
				}
			}
			cmd.Env = append(cmd.Env, helper+"=1")
			cmd.Env = append(cmd.Env, tc.env...)
			cmd.Stdin = strings.NewReader(tc.stdin)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("filter failed: %v\n%s", err, out)
			}
			if string(out) != tc.want {
				t.Errorf("got %q, want %q", out, tc.want)
			}
		})
	}
}
