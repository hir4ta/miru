package render

import (
	"os"
	"strings"
	"testing"
)

func TestToHTML(t *testing.T) {
	raw, err := os.ReadFile("../../sample.md")
	if err != nil {
		t.Fatal(err)
	}
	out, err := ToHTML("sample.md", string(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) < 1000 {
		t.Fatalf("html too small: %d bytes", len(out))
	}
	for _, want := range []string{
		"<!DOCTYPE html>",
		`class="markdown-body"`,
		"mumei-md sample",
		"<table>",
		"<code",
		"sample.md",
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("html missing %q", want)
		}
	}
	_ = os.WriteFile("/tmp/mumei-md-test.html", out, 0644)
	t.Logf("wrote /tmp/mumei-md-test.html (%d bytes)", len(out))
}
