package files

import "testing"

func TestShouldIgnore(t *testing.T) {
	for name, tt := range map[string]struct {
		path string
		want bool
	}{
		"git dir":              {path: ".git", want: true},
		"git nested":           {path: ".git/HEAD", want: true},
		"node_modules dir":     {path: "node_modules", want: true},
		"node_modules nested":  {path: "node_modules/foo/bar.js", want: true},
		"hidden file":          {path: ".env", want: true},
		"current dir prefix":   {path: ".", want: true},
		"regular file":         {path: "main.go", want: false},
		"nested regular":       {path: "src/foo.go", want: false},
		"file with dot inside": {path: "foo.bar", want: false},
	} {
		t.Run(name, func(t *testing.T) {
			if got := shouldIgnore(tt.path); got != tt.want {
				t.Errorf("shouldIgnore(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
