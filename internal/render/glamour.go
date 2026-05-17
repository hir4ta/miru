package render

import (
	"embed"
	"fmt"
	"os"

	"charm.land/glamour/v2"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
)

//go:embed assets/themes/*.json
var themeFS embed.FS

const DefaultTheme = "claude"

// Built-in Glamour standard styles available out-of-the-box.
var builtinStyles = map[string]bool{
	"dracula":     true,
	"tokyo-night": true,
	"dark":        true,
	"light":       true,
	"pink":        true,
	"ascii":       true,
	"notty":       true,
}

// AvailableThemes returns the list of theme names recognized by the renderer.
func AvailableThemes() []string {
	entries, _ := themeFS.ReadDir("assets/themes")
	var out []string
	for _, e := range entries {
		name := e.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			out = append(out, name[:len(name)-5])
		}
	}
	for name := range builtinStyles {
		out = append(out, name)
	}
	return out
}

type ANSI struct {
	renderer *glamour.TermRenderer
	width    int
	theme    string
}

func NewANSI(width int, theme string) (*ANSI, error) {
	r, err := newRenderer(width, theme)
	if err != nil {
		return nil, err
	}
	return &ANSI{renderer: r, width: width, theme: theme}, nil
}

func (a *ANSI) Render(markdown string) (string, error) {
	out, err := a.renderer.Render(markdown)
	if err != nil {
		return "", err
	}
	return PostProcessLists(markdown, out), nil
}

func (a *ANSI) Resize(width int) error {
	if width == a.width {
		return nil
	}
	r, err := newRenderer(width, a.theme)
	if err != nil {
		return err
	}
	a.renderer = r
	a.width = width
	return nil
}

func newRenderer(width int, theme string) (*glamour.TermRenderer, error) {
	if width < 20 {
		width = 80
	}
	if theme == "" {
		theme = DefaultTheme
	}
	// Glamour registers its chroma style under the hardcoded name "charm" but
	// skips re-registration on subsequent calls. Drop the cached entry so the
	// new theme's chroma palette is applied on the next render.
	delete(chromastyles.Registry, "charm")

	opts := []glamour.TermRendererOption{
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	}
	styleOpt, err := resolveTheme(theme)
	if err != nil {
		return nil, err
	}
	opts = append(opts, styleOpt)
	return glamour.NewTermRenderer(opts...)
}

func resolveTheme(theme string) (glamour.TermRendererOption, error) {
	if builtinStyles[theme] {
		return glamour.WithStandardStyle(theme), nil
	}
	bytes, err := themeFS.ReadFile("assets/themes/" + theme + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("unknown theme %q (available: %v)", theme, AvailableThemes())
		}
		return nil, fmt.Errorf("read theme %q: %w", theme, err)
	}
	return glamour.WithStylesFromJSONBytes(bytes), nil
}
