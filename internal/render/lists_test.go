package render

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestExtractOrderedMarkers(t *testing.T) {
	src := `## A
1. A
2. B
   1. nest 1
   2. nest 2
3. C

## B
1. X
2. Y

## C
1. P (mixed)
   - U1
   - U2
     1. P-inner-1
     2. P-inner-2
   - U3
2. Q
3. R
`
	got := extractOrderedMarkers(src)
	want := []string{
		"1.", "2.", "2.1", "2.2", "3.",
		"1.", "2.",
		"1.", "1.1", "1.2", "2.", "3.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestPostProcessLists_Visual(t *testing.T) {
	r, err := NewANSI(80, "claude")
	if err != nil {
		t.Fatal(err)
	}
	src := `## ul
- Item 1
- Item 2
  - Nested 2-1
  - Nested 2-2
    - Deep nested 2-2-1
- Item 3

## ol
1. First
2. Second
   1. Nested 2-1
   2. Nested 2-2
      1. Deep nested
3. Third

## mixed
1. Parent
   - Child
   - Child
     1. Grandchild
     2. Grandchild
   - Child
2. Parent
3. Parent
`
	out, err := r.Render(src)
	if err != nil {
		t.Fatal(err)
	}
	wantSubstrings := []string{
		"● Item 1",
		"○ Nested 2-1",
		"▪ Deep nested 2-2-1",
		"1. First",
		"2.1 Nested 2-1",
		"2.2 Nested 2-2",
		"2.2.1 Deep nested",
		"3. Third",
		"1. Parent",
		"1.1 Grandchild",
		"1.2 Grandchild",
		"2. Parent",
		"3. Parent",
	}
	stripped := ansi.Strip(out)
	for _, want := range wantSubstrings {
		if !strings.Contains(stripped, want) {
			t.Errorf("output missing %q", want)
		}
	}
	for i, line := range strings.Split(out, "\n") {
		s := ansi.Strip(line)
		if strings.TrimSpace(s) == "" {
			continue
		}
		t.Logf("[%2d] |%s|", i, s)
	}
}
