package nav

import (
	"os"
	"testing"

	"github.com/hir4ta/miru/internal/render"
)

func TestExtract(t *testing.T) {
	raw, err := os.ReadFile("../../sample.md")
	if err != nil {
		t.Fatal(err)
	}
	headings := Extract(string(raw))
	if len(headings) < 10 {
		t.Fatalf("expected many headings, got %d", len(headings))
	}
	if headings[0].Level != 1 || headings[0].Text != "miru sample" {
		t.Errorf("first heading mismatch: %+v", headings[0])
	}
	for i, h := range headings {
		t.Logf("[%2d] h%d: %s", i, h.Level, h.Text)
	}
}

func TestMapToLines(t *testing.T) {
	raw, err := os.ReadFile("../../sample.md")
	if err != nil {
		t.Fatal(err)
	}
	r, err := render.NewANSI(100, "claude")
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := r.Render(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	mapped := MapToLines(Extract(string(raw)), rendered)
	if len(mapped) < 5 {
		t.Fatalf("mapped too few: %d", len(mapped))
	}
	for i := 1; i < len(mapped); i++ {
		if mapped[i].Line <= mapped[i-1].Line {
			t.Errorf("line numbers not monotonic at %d: %d after %d", i, mapped[i].Line, mapped[i-1].Line)
		}
	}
	for i, h := range mapped {
		t.Logf("[%2d] line=%4d h%d: %s", i, h.Line, h.Level, h.Text)
	}
}

func TestPrevNext(t *testing.T) {
	hs := []Heading{
		{Level: 1, Text: "A", Line: 0},
		{Level: 2, Text: "B", Line: 10},
		{Level: 2, Text: "C", Line: 25},
		{Level: 1, Text: "D", Line: 50},
	}
	if n, ok := Next(hs, 5); !ok || n != 10 {
		t.Errorf("Next(5) = %d %v, want 10", n, ok)
	}
	if n, ok := Next(hs, 25); !ok || n != 50 {
		t.Errorf("Next(25) = %d %v, want 50", n, ok)
	}
	if _, ok := Next(hs, 60); ok {
		t.Error("Next(60) should fail")
	}
	if p, ok := Prev(hs, 12); !ok || p != 10 {
		t.Errorf("Prev(12) = %d %v, want 10", p, ok)
	}
	if p, ok := Prev(hs, 0); ok {
		t.Errorf("Prev(0) should fail, got %d", p)
	}
}
