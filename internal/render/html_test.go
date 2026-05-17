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
		"miru sample",
		"<table>",
		"<code",
		"sample.md",
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("html missing %q", want)
		}
	}
	_ = os.WriteFile("/tmp/miru-test.html", out, 0644)
	t.Logf("wrote /tmp/miru-test.html (%d bytes)", len(out))
}

func TestToHTML_MermaidGatesThirdPartyAssets(t *testing.T) {
	withMermaid := "# x\n\n```mermaid\nflowchart LR\n A --> B\n```\n"
	withoutMermaid := "# x\n\nplain markdown, no diagram.\n\n```go\nfunc main() {}\n```\n"

	out, err := ToHTML("x.md", withMermaid)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`look: "handDrawn"`,
		`fonts.googleapis.com/css2?family=Caveat`,
		`cdn.jsdelivr.net/npm/mermaid@11`,
		`font-family: "Caveat"`,
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("mermaid markdown: html missing %q", want)
		}
	}

	out, err = ToHTML("x.md", withoutMermaid)
	if err != nil {
		t.Fatal(err)
	}
	for _, dontWant := range []string{
		"mermaid",
		"fonts.googleapis.com",
		"fonts.gstatic.com",
		"Caveat",
	} {
		if strings.Contains(string(out), dontWant) {
			t.Errorf("non-mermaid markdown: html should not contain %q", dontWant)
		}
	}
}
