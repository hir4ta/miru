package tui

import (
	"fmt"
	"image/color"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	settingsBoxWidth  = 36
	settingsBoxHeight = 14
)

type themeItem struct{ name string }

func (i themeItem) FilterValue() string { return i.name }

type themeDelegate struct {
	accent color.Color
	muted  color.Color
}

func (d themeDelegate) Height() int                             { return 1 }
func (d themeDelegate) Spacing() int                            { return 0 }
func (d themeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d themeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(themeItem)
	if !ok {
		return
	}
	cursor := "  "
	style := lipgloss.NewStyle().Foreground(d.muted)
	if index == m.Index() {
		cursor = "> "
		style = lipgloss.NewStyle().Foreground(d.accent).Bold(true)
	}
	fmt.Fprint(w, style.Render(cursor+it.name))
}

type settingsModel struct {
	list list.Model
	open bool
}

func newSettings(themes []string, current string, accent, muted color.Color) settingsModel {
	items := make([]list.Item, len(themes))
	cursor := 0
	for i, t := range themes {
		items[i] = themeItem{name: t}
		if t == current {
			cursor = i
		}
	}
	l := list.New(items, themeDelegate{accent: accent, muted: muted}, settingsBoxWidth-4, settingsBoxHeight-4)
	l.Title = "Choose theme"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Foreground(accent).Bold(true)
	l.Select(cursor)
	return settingsModel{list: l}
}

func (m Model) settingsView() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.accent()).
		Padding(1, 2).
		Width(settingsBoxWidth).
		Height(settingsBoxHeight)

	hint := lipgloss.NewStyle().
		Foreground(m.muted()).
		Render("↑↓ select · ⏎ apply · esc cancel")

	content := m.settings.list.View() + "\n\n" + hint
	return box.Render(content)
}

func (m Model) overlaySettings(bg string) string {
	return overlayCentered(DimBackground(bg), m.settingsView(), m.winW, m.winH)
}

func overlayCentered(bg, box string, winW, winH int) string {
	boxW := lipgloss.Width(box)
	boxH := lipgloss.Height(box)
	x := max(0, (winW-boxW)/2)
	y := max(0, (winH-boxH)/2)
	return Overlay(bg, box, x, y)
}
