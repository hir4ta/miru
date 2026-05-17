package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestOverlay(t *testing.T) {
	bg := "AAAAAAAAAA\nBBBBBBBBBB\nCCCCCCCCCC\nDDDDDDDDDD\nEEEEEEEEEE"
	fg := "XX\nYY"
	out := Overlay(bg, fg, 2, 1)

	want := []string{
		"AAAAAAAAAA",
		"BBXXBBBBBB",
		"CCYYCCCCCC",
		"DDDDDDDDDD",
		"EEEEEEEEEE",
	}
	for i, line := range strings.Split(out, "\n") {
		got := ansi.Strip(line)
		if got != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got, want[i])
		}
	}
}

func TestOverlayWithANSI(t *testing.T) {
	bg := "\x1b[31mRED LINE\x1b[0m"
	fg := "X"
	out := Overlay(bg, fg, 3, 0)
	stripped := ansi.Strip(out)
	if stripped != "REDXLINE" {
		t.Errorf("got %q, want %q", stripped, "REDXLINE")
	}
}
