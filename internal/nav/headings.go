package nav

import (
	"bytes"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level int
	Text  string
	Line  int
}

func Extract(markdown string) []Heading {
	source := []byte(markdown)
	root := goldmark.New().Parser().Parse(text.NewReader(source))
	var out []Heading
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		out = append(out, Heading{
			Level: h.Level,
			Text:  headingText(h, source),
		})
		return ast.WalkSkipChildren, nil
	})
	return out
}

func headingText(h ast.Node, source []byte) string {
	var buf bytes.Buffer
	_ = ast.Walk(h, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := n.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}

func MapToLines(headings []Heading, rendered string) []Heading {
	lines := strings.Split(rendered, "\n")
	out := make([]Heading, 0, len(headings))
	cursor := 0
	for _, h := range headings {
		prefix := strings.Repeat("#", h.Level) + " "
		for ; cursor < len(lines); cursor++ {
			s := strings.TrimSpace(ansi.Strip(lines[cursor]))
			if strings.HasPrefix(s, prefix) {
				h.Line = cursor
				out = append(out, h)
				cursor++
				break
			}
		}
	}
	return out
}

func Prev(headings []Heading, currentLine int) (int, bool) {
	for i := len(headings) - 1; i >= 0; i-- {
		if headings[i].Line < currentLine {
			return headings[i].Line, true
		}
	}
	return 0, false
}

func Next(headings []Heading, currentLine int) (int, bool) {
	for _, h := range headings {
		if h.Line > currentLine {
			return h.Line, true
		}
	}
	return 0, false
}
