package tui

import "charm.land/lipgloss/v2"

func (m Model) helpView() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.accent()).
		Padding(1, 2)

	title := lipgloss.NewStyle().
		Foreground(m.accent()).
		Bold(true).
		Render("Help")

	body := m.help.FullHelpView(m.keys.FullHelp())

	hint := lipgloss.NewStyle().
		Foreground(m.muted()).
		Render("? close")

	return box.Render(title + "\n\n" + body + "\n\n" + hint)
}

func (m Model) overlayHelp(bg string) string {
	return overlayCentered(DimBackground(bg), m.helpView(), m.winW, m.winH)
}
