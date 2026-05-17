package render

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	docMargin   = 2
	levelIndent = 2
)

var bulletGlyphs = []string{"●", "○", "▪", "◦"}

var orderedItemRe = regexp.MustCompile(`^(\s*)(\d+)\. `)

// PostProcessLists rewrites Glamour's list output:
//   - Unordered: per-depth bullet glyph
//   - Ordered: hierarchical numbering from AST (1. / 1.1 / 1.2 / 2.)
func PostProcessLists(markdown, rendered string) string {
	markers := extractOrderedMarkers(markdown)

	lines := strings.Split(rendered, "\n")
	orderedQueue := 0
	for i, line := range lines {
		stripped := ansi.Strip(line)
		trimmed := strings.TrimLeft(stripped, " ")
		indent := len(stripped) - len(trimmed)
		depth := (indent - docMargin) / levelIndent
		if depth < 0 {
			depth = 0
		}

		// Unordered bullet: "• " is contiguous in raw (same style), simple replace works.
		if strings.HasPrefix(trimmed, "• ") {
			glyph := bulletGlyphs[depth%len(bulletGlyphs)]
			lines[i] = strings.Replace(line, "• ", glyph+" ", 1)
			continue
		}

		// Ordered "N. ": digits and ". " have different styles in raw, need ANSI-aware splice.
		m := orderedItemRe.FindStringSubmatchIndex(stripped)
		if m == nil {
			continue
		}
		if orderedQueue >= len(markers) {
			continue
		}
		visibleStart := m[4] // position of first digit in stripped
		visibleEnd := m[1]   // position right after ". "
		newMarker := markers[orderedQueue] + " "
		origVisibleMarker := stripped[visibleStart:visibleEnd]
		if newMarker != origVisibleMarker {
			lines[i] = spliceVisibleRange(line, visibleStart, visibleEnd, newMarker)
		}
		orderedQueue++
	}
	return strings.Join(lines, "\n")
}

// spliceVisibleRange replaces the visible characters in raw at positions
// [startVisible, endVisible) with replacement, preserving all surrounding
// ANSI escape sequences. ANSI escapes that fall WITHIN the visible range
// are also preserved (they are zero-width so this is safe).
func spliceVisibleRange(raw string, startVisible, endVisible int, replacement string) string {
	var sb strings.Builder
	sb.Grow(len(raw) + len(replacement))
	visible := 0
	i := 0
	inserted := false
	for i < len(raw) {
		if n := scanANSILen(raw[i:]); n > 0 {
			sb.WriteString(raw[i : i+n])
			i += n
			continue
		}
		if visible == startVisible && !inserted {
			sb.WriteString(replacement)
			inserted = true
		}
		if visible >= startVisible && visible < endVisible {
			visible++
			i++
			continue
		}
		sb.WriteByte(raw[i])
		visible++
		i++
	}
	if !inserted {
		sb.WriteString(replacement)
	}
	return sb.String()
}

// scanANSILen returns the byte length of an ANSI escape sequence starting at
// s[0], or 0 if s doesn't begin with an escape.
func scanANSILen(s string) int {
	if len(s) < 2 || s[0] != 0x1b {
		return 0
	}
	if s[1] == '[' {
		for i := 2; i < len(s); i++ {
			c := s[i]
			if c >= 0x40 && c <= 0x7e {
				return i + 1
			}
		}
		return len(s)
	}
	return 2
}

// extractOrderedMarkers walks the markdown AST and returns hierarchical
// markers for each ordered-list item in document order ("1.", "1.1", "1.2", "2.", ...).
func extractOrderedMarkers(markdown string) []string {
	root := goldmark.New().Parser().Parse(text.NewReader([]byte(markdown)))
	var (
		stack   []int
		markers []string
	)
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if list, ok := n.(*ast.List); ok && list.IsOrdered() {
			if entering {
				stack = append(stack, 0)
			} else if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
		if item, ok := n.(*ast.ListItem); ok && entering {
			if parent, ok := item.Parent().(*ast.List); ok && parent.IsOrdered() && len(stack) > 0 {
				stack[len(stack)-1]++
				parts := make([]string, len(stack))
				for j, c := range stack {
					parts[j] = strconv.Itoa(c)
				}
				var marker string
				if len(parts) == 1 {
					marker = parts[0] + "."
				} else {
					marker = strings.Join(parts, ".")
				}
				markers = append(markers, marker)
			}
		}
		return ast.WalkContinue, nil
	})
	return markers
}
