package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/hir4ta/mumei-md/internal/render"
)

type browserOpenedMsg struct{ err error }

func openInBrowser(filename, raw string) tea.Cmd {
	return func() tea.Msg {
		err := render.OpenInBrowser(filename, raw)
		return browserOpenedMsg{err: err}
	}
}
