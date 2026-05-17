package render

import (
	"encoding/json"

	"charm.land/glamour/v2/styles"
)

type themeColors struct {
	Heading struct {
		Color string `json:"color"`
	} `json:"heading"`
	BlockQuote struct {
		Color string `json:"color"`
	} `json:"block_quote"`
}

// AccentColor returns the primary accent (heading) color for the theme.
// Used for UI chrome — borders, filename label, etc.
func AccentColor(theme string) string {
	if theme == "" {
		theme = DefaultTheme
	}
	if c := lookupCustomThemeColor(theme, fieldHeading); c != "" {
		return c
	}
	if cfg, ok := styles.DefaultStyles[theme]; ok && cfg.Heading.Color != nil {
		return *cfg.Heading.Color
	}
	return "#d97757"
}

// MutedColor returns a dim secondary color for the theme.
// Used for the scroll percentage, secondary text.
func MutedColor(theme string) string {
	if theme == "" {
		theme = DefaultTheme
	}
	if c := lookupCustomThemeColor(theme, fieldBlockQuote); c != "" {
		return c
	}
	if cfg, ok := styles.DefaultStyles[theme]; ok && cfg.BlockQuote.Color != nil {
		return *cfg.BlockQuote.Color
	}
	return "#888888"
}

type colorField int

const (
	fieldHeading colorField = iota
	fieldBlockQuote
)

func lookupCustomThemeColor(theme string, field colorField) string {
	data, err := themeFS.ReadFile("assets/themes/" + theme + ".json")
	if err != nil {
		return ""
	}
	var t themeColors
	if err := json.Unmarshal(data, &t); err != nil {
		return ""
	}
	switch field {
	case fieldHeading:
		return t.Heading.Color
	case fieldBlockQuote:
		return t.BlockQuote.Color
	}
	return ""
}
