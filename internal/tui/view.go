package tui

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/hir4ta/miru/internal/render"
)

func (m Model) accent() color.Color {
	return lipgloss.Color(render.AccentColor(m.theme))
}

func (m Model) muted() color.Color {
	return lipgloss.Color(render.MutedColor(m.theme))
}

func (m Model) titleStyle() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().
		BorderStyle(b).
		BorderForeground(m.accent()).
		Foreground(m.accent()).
		Bold(true).
		Padding(0, 1)
}

func (m Model) infoStyle() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	return lipgloss.NewStyle().
		BorderStyle(b).
		BorderForeground(m.accent()).
		Foreground(m.muted()).
		Padding(0, 1)
}

func (m Model) dividerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(m.accent())
}

var helpBarStyle = lipgloss.NewStyle().Padding(0, 1)

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	if !m.ready {
		v.SetContent("\n  Initializing...")
		return v
	}

	body := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		helpBarStyle.Render(m.help.View(m.keys)),
	)

	switch {
	case m.settings.open:
		body = m.overlaySettings(body)
	case m.helpOpen:
		body = m.overlayHelp(body)
	}

	v.SetContent(body)
	return v
}

func (m Model) headerView() string {
	title := m.titleStyle().Render(filepath.Base(m.filename))
	width := m.viewport.Width()
	pad := max(0, width-lipgloss.Width(title))
	line := m.dividerStyle().Render(strings.Repeat("─", pad))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	pct := m.viewport.ScrollPercent() * 100
	info := m.infoStyle().Render(fmt.Sprintf("%3.0f%% · %s", pct, m.theme))
	width := m.viewport.Width()
	pad := max(0, width-lipgloss.Width(info))
	line := m.dividerStyle().Render(strings.Repeat("─", pad))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
