package render

import (
	"strings"
	"testing"
)

// Verify that switching themes actually changes the chroma code-block colors.
// Regression for Glamour's hardcoded "charm" registry name — without the
// cache-clear in newRenderer, all themes would share the first one's palette.
func TestChromaSwitchesPerTheme(t *testing.T) {
	src := "```go\nfunc main() { fmt.Println(\"hi\") }\n```\n"

	outs := map[string]string{}
	for _, theme := range []string{"claude", "gruvbox", "nord"} {
		r, err := NewANSI(80, theme)
		if err != nil {
			t.Fatal(err)
		}
		o, err := r.Render(src)
		if err != nil {
			t.Fatal(err)
		}
		outs[theme] = o
	}

	for a, oa := range outs {
		for b, ob := range outs {
			if a < b && oa == ob {
				t.Errorf("theme %q and %q produced identical output", a, b)
			}
		}
	}

	// Sanity: each output should contain SOME foreground color escape
	// (proves chroma applied a style at all, not just empty render).
	for theme, body := range outs {
		if !strings.Contains(body, "\x1b[38;5;") {
			t.Errorf("theme %q has no foreground 256-color escape", theme)
		}
	}
}
