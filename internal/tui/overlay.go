package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// dimSeq styles the background text behind a floating panel: dim
// attribute + bright-black (grey). Strips all colour from the
// underlying content so the floater reads as foreground.
const dimSeq = "\x1b[2;90m"

// DimBackground re-renders bg with all original styling stripped and
// replaced by a uniform dim grey. The terminal can't truly blur, so
// monochrome-dim is the closest readable approximation.
func DimBackground(bg string) string {
	lines := strings.Split(bg, "\n")
	out := make([]string, len(lines))
	for i, line := range lines {
		stripped := ansi.Strip(line)
		if stripped == "" {
			out[i] = ""
			continue
		}
		out[i] = dimSeq + stripped + "\x1b[0m"
	}
	return strings.Join(out, "\n")
}

// Overlay draws fg onto bg starting at terminal column x, row y.
// Both strings may contain ANSI escape sequences; visible columns are
// counted using the standard width rules.
func Overlay(bg, fg string, x, y int) string {
	bgLines := strings.Split(bg, "\n")
	for len(bgLines) < y+lipgloss.Height(fg) {
		bgLines = append(bgLines, "")
	}
	for i, fgLine := range strings.Split(fg, "\n") {
		bgLines[y+i] = overlayLine(bgLines[y+i], x, fgLine)
	}
	return strings.Join(bgLines, "\n")
}

func overlayLine(bg string, x int, fg string) string {
	fgW := lipgloss.Width(fg)
	bgW := lipgloss.Width(bg)

	left := ansi.Cut(bg, 0, x)
	if w := lipgloss.Width(left); w < x {
		left += strings.Repeat(" ", x-w)
	}
	right := ""
	if x+fgW < bgW {
		right = ansi.Cut(bg, x+fgW, bgW)
	}
	const reset = "\x1b[0m"
	return left + reset + fg + reset + right
}
