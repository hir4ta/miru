package render

import (
	"slices"
	"strings"
	"testing"
)

func TestAvailableThemes(t *testing.T) {
	themes := AvailableThemes()
	want := []string{"claude", "gruvbox", "everforest", "nord", "dracula", "tokyo-night"}
	for _, w := range want {
		if !slices.Contains(themes, w) {
			t.Errorf("missing theme %q in %v", w, themes)
		}
	}
}

func TestNewANSI_AllThemesLoadable(t *testing.T) {
	for _, theme := range AvailableThemes() {
		t.Run(theme, func(t *testing.T) {
			r, err := NewANSI(80, theme)
			if err != nil {
				t.Fatalf("NewANSI(%q): %v", theme, err)
			}
			out, err := r.Render("# title\n\nbody\n")
			if err != nil {
				t.Fatalf("Render(%q): %v", theme, err)
			}
			if !strings.Contains(out, "title") {
				t.Errorf("theme %q rendered no title", theme)
			}
		})
	}
}

func TestNewANSI_UnknownTheme(t *testing.T) {
	_, err := NewANSI(80, "does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown theme")
	}
	if !strings.Contains(err.Error(), "unknown theme") {
		t.Errorf("expected 'unknown theme' in error, got: %v", err)
	}
}
